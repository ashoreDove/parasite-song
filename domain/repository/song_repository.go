package repository

import (
	"github.com/ashoreDove/parasite-song/domain/model"
	"github.com/jinzhu/gorm"
)

type ISongRepository interface {
	InitTable() error
	FindSongByID(int64) (*model.Song, error)
	CreateSong(*model.Song) (int64, error)
	DeleteSongByID(int64) error
	UpdateSong(*model.Song) error
	FindAll() ([]model.Song, error)
	FindSongByName(string) ([]model.Song, error)
}

//创建songRepository
func NewSongRepository(db *gorm.DB) ISongRepository {
	return &SongRepository{mysqlDb: db}
}

type SongRepository struct {
	mysqlDb *gorm.DB
}

func (u *SongRepository) FindSongByName(name string) (songList []model.Song, err error) {
	return songList, u.mysqlDb.Where("song_name like ?", "%"+name+"%").Find(&songList).Error
}

//初始化表
func (u *SongRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Song{}).Error
}

//根据ID查找Song信息
func (u *SongRepository) FindSongByID(songID int64) (song *model.Song, err error) {
	song = &model.Song{}
	return song, u.mysqlDb.First(song, songID).Error
}

//创建Song信息
func (u *SongRepository) CreateSong(song *model.Song) (int64, error) {
	return song.SongId, u.mysqlDb.Create(song).Error
}

//根据ID删除Song信息
func (u *SongRepository) DeleteSongByID(songID int64) error {
	return u.mysqlDb.Where("id = ?", songID).Delete(&model.Song{}).Error
}

//更新Song信息
func (u *SongRepository) UpdateSong(song *model.Song) error {
	return u.mysqlDb.Model(song).Update(song).Error
}

//获取结果集
func (u *SongRepository) FindAll() (songAll []model.Song, err error) {
	return songAll, u.mysqlDb.Find(&songAll).Error
}
