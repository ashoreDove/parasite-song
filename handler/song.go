package handler

import (
	"context"
	"github.com/ashoreDove/common"
	"github.com/ashoreDove/parasite-song/domain/service"
	song "github.com/ashoreDove/parasite-song/proto/song"
	"github.com/ashoreDove/parasite-song/utils"
	"github.com/jinzhu/gorm"
	"net/http"
)

type ISongService interface {
	Search(context.Context, *song.SearchRequest, *song.SearchResponse) error
	GetSongInfo(context.Context, *song.SongIdRequest, *song.SongResponse) error
	GetSongList(context.Context, *song.ListRequest, *song.SongListResponse) error
}
type Song struct {
	SongDataService service.ISongDataService
	client          *http.Client
}

func (s Song) Search(ctx context.Context, req *song.SearchRequest, resp *song.SearchResponse) error {
	songList, err := s.SongDataService.FindSongByName(req.Keyword)
	if len(songList) == 0 {
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

func (s Song) GetSongInfo(ctx context.Context, request *song.SongIdRequest, response *song.SongResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s Song) GetSongList(ctx context.Context, request *song.ListRequest, response *song.SongListResponse) error {
	//TODO implement me
	panic("implement me")
}

func NewSongService(db *gorm.DB) ISongService {
	return &Song{SongDataService: service.NewSongDataService(db), client: &http.Client{}}
}
