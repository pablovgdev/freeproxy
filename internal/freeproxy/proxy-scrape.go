package freeproxy

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ProxyScrape struct{}

type proxyScrapeResponse struct {
	Proxies []proxyScrapeProxy `json:"proxies"`
}

type proxyScrapeProxy struct {
	IP             string `json:"ip"`
	Port           int    `json:"port"`
	AnonymityLevel string `json:"anonymity"`
	Protocol       string `json:"protocol"`
	SSL            bool   `json:"ssl"`
	IPData         IPData `json:"ip_data"`
}

type IPData struct {
	CountryCode string `json:"countryCode"`
}

func (s *ProxyScrape) GetProxies(params GetProxiesParams) []Proxy {
	baseURL := "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&timeout=10000&proxy_format=ipport&format=json"
	url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
		return []Proxy{}
	}

	query := url.Query()

	if params.CountryCode != "" {
		query.Set("country", params.CountryCode)
	}

	if params.AnonimityLevel != 0 {
		switch params.AnonimityLevel {
		case Transparent:
			query.Set("anonymity", "transparent")
		case Anonymous:
			query.Set("anonymity", "anonymous")
		case Elite:
			query.Set("anonymity", "elite")
		default:
			query.Set("anonymity", "transparent")
		}
	}

	if params.Protocol != 0 {
		switch params.Protocol {
		case HTTP:
		case HTTPS:
			query.Set("protocol", "http")
		case Socks4:
			query.Set("protocol", "socks4")
		case Socks5:
			query.Set("protocol", "socks5")
		default:
			query.Set("protocol", "http")
		}
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

	var proxyScrapeResp proxyScrapeResponse
	err = json.Unmarshal(body, &proxyScrapeResp)
	if err != nil {
		log.Fatal(err)
		return proxies
	}

	for _, proxy := range proxyScrapeResp.Proxies {
		proxies = append(proxies, s.toProxy(proxy))
	}

	return proxies
}

func (s *ProxyScrape) toProxy(p proxyScrapeProxy) Proxy {
	var anonymityLevel ProxyAnonimityLevel

	switch p.AnonymityLevel {
	case "transparent":
		anonymityLevel = Transparent
	case "anonymous":
		anonymityLevel = Anonymous
	case "elite":
		anonymityLevel = Elite
	default:
		anonymityLevel = Transparent
	}

	var protocol ProxyProtocol

	switch p.Protocol {
	case "http":
		if p.SSL {
			protocol = HTTPS
		} else {
			protocol = HTTP
		}
		protocol = HTTPS
	case "socks4":
		protocol = Socks4
	case "socks5":
		protocol = Socks5
	default:
		protocol = HTTP
	}

	return *NewProxy(p.IP, p.Port, anonymityLevel, protocol, p.IPData.CountryCode, "ProxyScrape")
}
