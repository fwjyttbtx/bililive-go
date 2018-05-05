package api

import (
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/lib/http"
	"github.com/tidwall/gjson"
	"net/url"
	"regexp"
	"strings"
)

const (
	longzhuMobileUrl  = "http://m.longzhu.com/"
	longzhuRoomApiUrl = "http://liveapi.plu.cn/liveapp/roomstatus"
	longzhuLiveApiUrl = "http://livestream.plu.cn/live/getlivePlayurl"
)

type LongzhuLive struct {
	Url    *url.URL
	realId string
}

func (l *LongzhuLive) parseRealId() error {
	dom, err := http.Get(fmt.Sprintf("%s%s", longzhuMobileUrl, strings.Split(l.Url.Path, "/")[1]), nil, nil)
	if err != nil {
		return err
	}
	reg := regexp.MustCompile(`var\s*roomId\s*=\s*(\d+);`)
	realIds := reg.FindStringSubmatch(string(dom))
	if realIds == nil {
		return &RoomNotExistsError{l.Url}
	}
	l.realId = realIds[1]
	return nil
}

func (l *LongzhuLive) GetRoom() (*Info, error) {
	if l.realId == "" {
		if err := l.parseRealId(); err != nil {
			return nil, err
		}
	}
	body, err := http.Get(longzhuRoomApiUrl, map[string]string{"roomId": l.realId}, nil)
	if err != nil {
		return nil, err
	}
	info := &Info{
		Live:     l,
		Url:      l.Url,
		HostName: gjson.GetBytes(body, "userName").String(),
		RoomName: gjson.GetBytes(body, "title").String(),
		Status:   len(gjson.GetBytes(body, "streamUri").String()) > 4,
	}
	return info, nil

}

func (l *LongzhuLive) GetUrls() ([]*url.URL, error) {
	if l.realId == "" {
		if err := l.parseRealId(); err != nil {
			return nil, err
		}
	}
	body, err := http.Get(longzhuLiveApiUrl, map[string]string{"roomId": l.realId}, nil)
	if err != nil {
		return nil, err
	}
	urls := gjson.GetBytes(body, "playLines.0.urls.#.securityUrl").Array()
	us := make([]*url.URL, 0, 4)
	for _, u := range urls {
		url_, _ := url.Parse(u.String())
		us = append(us, url_)
	}
	return us, nil
}
