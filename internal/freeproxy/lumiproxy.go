package freeproxy

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type LumiproxyProxyService struct{}

type lumiproxyResponse struct {
	RequestID string        `json:"request_id"`
	Code      int           `json:"code"`
	Msg       string        `json:"msg"`
	Data      lumiproxyData `json:"data"`
}

type lumiproxyData struct {
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	PageCount  int              `json:"page_count"`
	TotalCount int              `json:"total_count"`
	List       []lumiproxyProxy `json:"list"`
}

type lumiproxyProxy struct {
	IP             string `json:"ip"`
	Port           int    `json:"port"`
	AnonymityLevel int    `json:"anonymity"`
	Protocol       int    `json:"protocol"`
	Speed          int    `json:"speed"`
	Uptime         int    `json:"uptime"`
	Latency        int    `json:"latency"`
	GooglePassed   int    `json:"google_passed"`
	CountryCode    string `json:"country_code"`
	Valid          bool
	ResponseTime   int
}

func (s *LumiproxyProxyService) GetProxies(params GetProxiesParams) []Proxy {
	baseURL := "https://api.lumiproxy.com/web_v1/free-proxy/list"
	url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
		return []Proxy{}
	}

	query := url.Query()

	if params.Language != "" {
		query.Set("language", params.Language)
	} else {
		query.Set("language", "en-us")
	}

	if params.PageSize != 0 {
		query.Set("page_size", strconv.Itoa(params.PageSize))
	} else {
		query.Set("page_size", "1000")
	}

	if params.Page != 0 {
		query.Set("page", strconv.Itoa(params.Page))
	} else {
		query.Set("page", "1")
	}

	if params.CountryCode != "" {
		query.Set("country_code", params.CountryCode)
	}

	if params.AnonimityLevel != 0 {
		query.Set("anonymity", strconv.Itoa(int(params.AnonimityLevel)))
	}

	if params.Protocol != 0 {
		query.Set("protocol", strconv.Itoa(int(params.Protocol)))
	}

	if params.Speed != 0 {
		query.Set("speed", strconv.Itoa(int(params.Speed)))
	}

	if params.Uptime != 0 {
		query.Set("uptime", strconv.Itoa(params.Uptime))
	}

	if params.GooglePassed {
		query.Set("google_passed", "1")
	}

	var proxies []Proxy

	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
		return proxies
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return proxies
	}

	var lumiproxyResp lumiproxyResponse
	err = json.Unmarshal(body, &lumiproxyResp)
	if err != nil {
		log.Fatal(err)
		return proxies
	}

	lumiproxyProxies := lumiproxyResp.Data.List

	for _, lumiproxyProxy := range lumiproxyProxies {
		proxies = append(proxies, s.toProxy(lumiproxyProxy))
	}

	return proxies
}

func (s *LumiproxyProxyService) toProxy(lumiproxyProxy lumiproxyProxy) Proxy {
	proxy := Proxy{
		IP:             lumiproxyProxy.IP,
		Port:           lumiproxyProxy.Port,
		AnonymityLevel: ProxyAnonimityLevel(lumiproxyProxy.AnonymityLevel),
		Protocol:       ProxyProtocol(lumiproxyProxy.Protocol),
		Speed:          ProxySpeed(lumiproxyProxy.Speed),
		Uptime:         lumiproxyProxy.Uptime,
		Latency:        lumiproxyProxy.Latency,
		GooglePassed:   lumiproxyProxy.GooglePassed,
		CountryCode:    lumiproxyProxy.CountryCode,
		Valid:          lumiproxyProxy.Valid,
		ResponseTime:   lumiproxyProxy.ResponseTime,
	}

	switch lumiproxyProxy.AnonymityLevel {
	case 0:
		proxy.AnonymityLevel = Transparent
	case 1:
		proxy.AnonymityLevel = Anonymous
	case 2:
		proxy.AnonymityLevel = Elite
	default:
		proxy.AnonymityLevel = Transparent
	}

	switch lumiproxyProxy.Protocol {
	case 1:
		proxy.Protocol = HTTP
	case 2:
		proxy.Protocol = HTTPS
	case 4:
		proxy.Protocol = Socks4
	case 8:
		proxy.Protocol = Socks5
	default:
		proxy.Protocol = HTTP
	}

	switch lumiproxyProxy.Speed {
	case 0:
		proxy.Speed = Slow
	case 1:
		proxy.Speed = Medium
	case 2:
		proxy.Speed = Fast
	default:
		proxy.Speed = Slow
	}

	return proxy
}
