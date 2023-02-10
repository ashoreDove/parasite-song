package model

import (
	"strconv"
	"strings"
)

type Song struct {
	SongId   int64  `gorm:"primary_key;not_null;auto_increment" json:"song_id"`
	SongName string `gorm:"not_null" json:"song_name"`
	Artist   string `gorm:"not_null" json:"artist"`
	Total    string `gorm:"not_null" json:"total"`
	IsTmp    int64  `gorm:"not_null" json:"is_tmp"`
}

type SongModel struct {
	Sid       int64  `json:"song_id"`
	SongName  string `json:"song_name"`
	Artist    string `json:"artist"`
	TotalTime string `json:"song_time_minutes"`
	Url       string `json:"url"`
}

func SongToModel(s SongModel) *Song {
	isNoTmp := strings.Contains(s.Url, "192.168.0.106")
	isTmp := 1
	if isNoTmp {
		isTmp = 0
	}
	return &Song{
		SongId:   s.Sid,
		SongName: s.SongName,
		Artist:   s.Artist,
		Total:    s.TotalTime,
		IsTmp:    int64(isTmp),
	}
}
func ModelToSong(s Song) *SongModel {
	url := ""
	if s.IsTmp == 0 {
		//todo 这边暂时这样处理
		url = "http://192.168.0.106.1/" + strconv.FormatInt(s.SongId, 10) + ".mp3"
	}
	return &SongModel{
		Sid:       s.SongId,
		SongName:  s.SongName,
		Artist:    s.Artist,
		TotalTime: s.Total,
		Url:       url,
	}
}
