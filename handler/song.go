package handler

import (
	"bytes"
	"context"
	"github.com/ashoreDove/common"
	"github.com/ashoreDove/parasite-song/domain/model"
	"github.com/ashoreDove/parasite-song/domain/service"
	song "github.com/ashoreDove/parasite-song/proto/song"
	"github.com/ashoreDove/parasite-song/utils"
	"github.com/jinzhu/gorm"
	"github.com/jlaffaye/ftp"
	log "github.com/micro/go-micro/v2/logger"
	"net/http"
	"strconv"
)

type ISongService interface {
	Search(context.Context, *song.SearchRequest, *song.SearchResponse) error
	GetSongInfo(context.Context, *song.SongIdRequest, *song.SongResponse) error
	GetSongList(context.Context, *song.ListRequest, *song.SongListResponse) error
}
type Song struct {
	SongDataService service.ISongDataService
	client          *http.Client
	ftpConn         *ftp.ServerConn
}

func (s Song) Search(ctx context.Context, req *song.SearchRequest, resp *song.SearchResponse) error {
	songList, err := s.SongDataService.FindSongByName(req.Keyword)
	if len(songList) < 5 {
		//抓包
		songList, err = utils.Search(s.client, req.Keyword)
		if err != nil {
			return err
		}
		//开线程写入数据库
		go func() {
			err := s.SongDataService.AddSongList(songList)
			if err != nil {
				panic(err)
			}
		}()
	}
	for _, v := range songList {
		s := &song.SongInfo{}
		if err := common.SwapTo(v, s); err != nil {
			return err
		}
		resp.SongList = append(resp.SongList, s)
	}

	return err
}

func (s Song) GetSongInfo(ctx context.Context, req *song.SongIdRequest, resp *song.SongResponse) error {
	songInfo, err := s.SongDataService.FindSongByID(req.SongId)
	if err != nil {
		return err
	}
	resp.SongInfo = &song.SongInfo{
		SongId:          songInfo.SongId,
		SongName:        songInfo.SongName,
		SongTimeMinutes: songInfo.Total,
		Artist:          songInfo.Artist,
	}
	if songInfo.IsTmp == 1 {
		//抓包获取Url和内容
		sUrl, err := utils.GetSongById(s.client, req.SongId)
		if err != nil {
			resp.IsSuccess = false
			return err
		}
		resp.IsSuccess = true
		resp.SongInfo.SongUrl = sUrl
		//存入ftp服务器
		//更新数据库
		go func() {
			err := s.Upload(sUrl, songInfo)
			for err != nil {
				err = s.Upload(sUrl, songInfo)
			}
		}()
		return nil
	}
	resp.IsSuccess = true
	resp.SongInfo.SongUrl = "http://192.168.0.106/" + strconv.FormatInt(songInfo.SongId, 10) + ".mp3"
	return nil
}
func (s Song) Upload(s_url string, s_model *model.Song) error {
	byt, err := utils.GetSongContent(s.client, s_url)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(*byt)
	err = s.ftpConn.Stor(strconv.FormatInt(s_model.SongId, 10)+".mp3", reader)
	if err != nil {
		return err
	}
	log.Info("创建文件成功")
	s_model.IsTmp = 0
	err = s.SongDataService.UpdateSong(s_model)
	if err != nil {
		return err
	}
	return nil
}

func (s Song) GetSongList(ctx context.Context, request *song.ListRequest, response *song.SongListResponse) error {
	//TODO implement me
	panic("implement me")
}

func NewSongService(db *gorm.DB, conn *ftp.ServerConn) ISongService {
	return &Song{SongDataService: service.NewSongDataService(db), client: &http.Client{}, ftpConn: conn}
}
