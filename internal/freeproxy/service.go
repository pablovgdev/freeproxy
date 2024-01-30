package freeproxy

import (
	"sync"
)

type Service struct {
	Lumiproxy   ProxyService
	ProxyScrape ProxyService
}

func New() *Service {
	return &Service{
		Lumiproxy:   &Lumiproxy{},
		ProxyScrape: &ProxyScrape{},
	}
}

func (s *Service) GetProxies(params GetProxiesParams) []Proxy {
	proxies := []Proxy{}

	lumiproxyProxies := s.Lumiproxy.GetProxies(params)
	proxies = append(proxies, lumiproxyProxies...)

	proxyScrapeProxies := s.ProxyScrape.GetProxies(params)
	proxies = append(proxies, proxyScrapeProxies...)

	return proxies
}

func (s *Service) ValidateProxies(freeProxyList *[]Proxy) {
	wg := sync.WaitGroup{}

	for i := range *freeProxyList {
		wg.Add(1)
		proxy := &(*freeProxyList)[i]

		go func(i int, proxy *Proxy) {
			defer wg.Done()
			proxy.Validate()
		}(i, proxy)
	}

	wg.Wait()
}
