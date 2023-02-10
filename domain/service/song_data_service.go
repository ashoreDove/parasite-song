package service

import (
	"github.com/ashoreDove/parasite-song/domain/model"
	"github.com/ashoreDove/parasite-song/domain/repository"
	"github.com/jinzhu/gorm"
)

type ISongDataService interface {
	AddSong(model1 *model.Song) (int64, error)
	UpdateSong(*model.Song) error
	FindSongByID(int64) (*model.Song, error)
	FindSongByName(string) ([]model.SongModel, error)
	AddSongList([]model.SongModel) error
}

//创建
func NewSongDataService(db *gorm.DB) ISongDataService {
	return &SongDataService{repository.NewSongRepository(db)}
}

type SongDataService struct {
	SongRepository repository.ISongRepository
}

func (u *SongDataService) AddSongList(list []model.SongModel) error {
	for _, v := range list {
		_, err := u.SongRepository.FindSongByID(v.Sid)
		if err == nil {
			//找到了
			continue
		}
		_, err = u.SongRepository.CreateSong(model.SongToModel(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *SongDataService) FindSongByName(name string) (songList []model.SongModel, err error) {
	songModelList, err := u.SongRepository.FindSongByName(name)
	if err != nil {
		return nil, err
	}
	for _, value := range songModelList {
		songList = append(songList, *model.ModelToSong(value))
	}
	return songList, err
}

//插入
func (u *SongDataService) AddSong(song *model.Song) (int64, error) {
	return u.SongRepository.CreateSong(song)
}

//更新
func (u *SongDataService) UpdateSong(song *model.Song) error {
	return u.SongRepository.UpdateSong(song)
}

//查找
func (u *SongDataService) FindSongByID(songID int64) (*model.Song, error) {
	return u.SongRepository.FindSongByID(songID)
}
