package utils

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ashoreDove/parasite-song/domain/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	searchUrl  = "https://www.kuwo.cn/api/www/search/searchMusicBykeyWord?rn=20&httpsStatus=1&reqId=51a8bf60-8396-11ed-888c-4197747a293a&key="
	songSrcUrl = "http://www.kuwo.cn/api/v1/www/music/playUrl?type=convert_url3&br=320kmp3&mid="
	lrcUrl     = "http://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId="
)

var header = map[string]string{
	"Referer":    "https://www.kuwo.cn/search/list?key=",
	"cookie":     "kw_token=SA4RWNUIKT8",
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.54",
	"csrf":       "SA4RWNUIKT8",
}

func Search(client *http.Client, key string) (sList []model.SongModel, err error) {
	keyUrl := url.QueryEscape(key)
	req, err := http.NewRequest("GET", searchUrl+keyUrl, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	req = req.WithContext(ctx)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	defer cancelFunc()
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var respBody map[string]interface{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return nil, err
	}
	data := respBody["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	for i, _ := range list {
		item := list[i].(map[string]interface{})
		//pic
		song := model.SongModel{
			Sid:       int64(item["rid"].(float64)),
			SongName:  item["name"].(string),
			Artist:    item["artist"].(string),
			TotalTime: item["songTimeMinutes"].(string),
			Url:       "",
		}
		sList = append(sList, song)
	}
	return sList, nil
}

func GetSongById(client *http.Client, id int64) (string, error) {
	url := songSrcUrl + strconv.FormatInt(id, 10)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	req = req.WithContext(ctx)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	defer cancelFunc()
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var respBody map[string]interface{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return "", err
	}
	code := respBody["code"].(float64)
	if code == 200 {
		data := respBody["data"].(map[string]interface{})
		return data["url"].(string), nil
	}
	return "", errors.New("????????????url??????")
}

func GetSongContent(client *http.Client, s_url string) (*[]byte, error) {
	req, err := http.NewRequest("GET", s_url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	for k, v := range header {
		req.Header.Add(k, v)
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	defer cancelFunc()
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return &body, nil
}
