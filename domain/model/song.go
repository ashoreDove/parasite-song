package model

import (
	"strconv"
	"strings"
)

type SongGormModel struct {
	SongId   int64  `gorm:"primary_key;not_null;auto_increment" json:"song_id"`
	SongName string `gorm:"not_null" json:"song_name"`
	Artist   string `gorm:"not_null" json:"artist"`
	Total    string `gorm:"not_null" json:"total"`
	IsTmpUrl bool   `gorm:"not_null" json:"is_tmp"`
}

type Song struct {
	Sid       int64  `json:"song_id"`
	SongName  string `json:"song_name"`
	Artist    string `json:"artist"`
	TotalTime string `json:"song_time_minutes"`
	Url       string `json:"url"`
}

func SongToModel(s Song) *SongGormModel {
	isNoTmp := strings.Contains(s.Url, "172.19.96.1")
	return &SongGormModel{
		SongId:   s.Sid,
		SongName: s.SongName,
		Artist:   s.Artist,
		Total:    s.TotalTime,
		IsTmpUrl: !isNoTmp,
	}
}
func ModelToSong(s SongGormModel) *Song {
	url := ""
	if !s.IsTmpUrl {
		//todo 这边暂时这样处理
		url = "http://172.19.96.1/" + strconv.FormatInt(s.SongId, 10) + ".mp3"
	}
	return &Song{
		Sid:       s.SongId,
		SongName:  s.SongName,
		Artist:    s.Artist,
		TotalTime: s.Total,
		Url:       url,
	}
}