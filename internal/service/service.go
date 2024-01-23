package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/pablovgdev/freeproxy/internal/models"
)

type Service struct{}

type GetProxiesParams struct {
	Language       string                     `json:"language"`
	PageSize       int                        `json:"page_size"`
	Page           int                        `json:"page"`
	CountryCode    string                     `json:"country_code"`
	AnonimityLevel models.ProxyAnonimityLevel `json:"anonymity"`
	Protocol       models.ProxyProtocol       `json:"protocol"`
	Speed          models.ProxySpeed          `json:"speed"`
	Uptime         int                        `json:"uptime"`
	GooglePassed   bool                       `json:"google_passed"`
}

type GetProxiesResponse struct {
	RequestID string         `json:"request_id"`
	Code      int            `json:"code"`
	Msg       string         `json:"msg"`
	Data      GetProxiesData `json:"data"`
}

type GetProxiesData struct {
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	PageCount  int            `json:"page_count"`
	TotalCount int            `json:"total_count"`
	List       []models.Proxy `json:"list"`
}

func (s *Service) GetProxies(params GetProxiesParams) []models.Proxy {
	baseURL := "https://api.lumiproxy.com/web_v1/free-proxy/list"
	url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
		return []models.Proxy{}
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

	var proxyList []models.Proxy

	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
		return proxyList
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return proxyList
	}

	var getProxiesResponse GetProxiesResponse
	err = json.Unmarshal(body, &getProxiesResponse)
	if err != nil {
		log.Fatal(err)
		return proxyList
	}

	return getProxiesResponse.Data.List
}

func (s *Service) ValidateProxies(freeProxyList *[]models.Proxy) {
	wg := sync.WaitGroup{}

	for i := range *freeProxyList {
		wg.Add(1)
		proxy := &(*freeProxyList)[i]

		go func(i int, proxy *models.Proxy) {
			defer wg.Done()
			proxy.Validate()
		}(i, proxy)
	}

	wg.Wait()
}
