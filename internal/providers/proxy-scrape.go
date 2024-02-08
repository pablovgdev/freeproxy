package providers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/pablovgdev/freeproxy/internal/proxy"
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
	IPData         ipData `json:"ip_data"`
}

type ipData struct {
	CountryCode string `json:"countryCode"`
}

func (ps *ProxyScrape) GetProxies(params proxy.GetProxiesParams) []proxy.Proxy {
	baseURL := "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&timeout=10000&proxy_format=ipport&format=json"
	url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
		return []proxy.Proxy{}
	}

	query := url.Query()

	if params.CountryCode != "" {
		query.Set("country", params.CountryCode)
	}

	if params.AnonimityLevel != 0 {
		switch params.AnonimityLevel {
		case proxy.Transparent:
			query.Set("anonymity", "transparent")
		case proxy.Anonymous:
			query.Set("anonymity", "anonymous")
		case proxy.Elite:
			query.Set("anonymity", "elite")
		default:
			query.Set("anonymity", "transparent")
		}
	}

	if params.Protocol != 0 {
		switch params.Protocol {
		case proxy.HTTP:
		case proxy.HTTPS:
			query.Set("protocol", "http")
		case proxy.Socks4:
			query.Set("protocol", "socks4")
		case proxy.Socks5:
			query.Set("protocol", "socks5")
		default:
			query.Set("protocol", "http")
		}
	}

	var proxies []proxy.Proxy

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
		proxies = append(proxies, ps.toProxy(proxy))
	}

	return proxies
}

func (ps *ProxyScrape) toProxy(p proxyScrapeProxy) proxy.Proxy {
	var anonymityLevel proxy.ProxyAnonimityLevel

	switch p.AnonymityLevel {
	case "transparent":
		anonymityLevel = proxy.Transparent
	case "anonymous":
		anonymityLevel = proxy.Anonymous
	case "elite":
		anonymityLevel = proxy.Elite
	default:
		anonymityLevel = proxy.Transparent
	}

	var protocol proxy.ProxyProtocol

	switch p.Protocol {
	case "http":
		if p.SSL {
			protocol = proxy.HTTPS
		} else {
			protocol = proxy.HTTP
		}
		protocol = proxy.HTTPS
	case "socks4":
		protocol = proxy.Socks4
	case "socks5":
		protocol = proxy.Socks5
	default:
		protocol = proxy.HTTP
	}

	return *proxy.NewProxy(p.IP, p.Port, anonymityLevel, protocol, p.IPData.CountryCode, "ProxyScrape")
}
