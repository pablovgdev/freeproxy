package freeproxy

import (
	"sync"
)

type Service struct {
	Lumiproxy ProxyService
}

func New() *Service {
	return &Service{
		Lumiproxy: &LumiproxyProxyService{},
	}
}

func (s *Service) GetProxyList(params GetProxiesParams) []Proxy {
	proxies := []Proxy{}

	lumiproxyProxies := s.Lumiproxy.GetProxies(params)

	proxies = append(proxies, lumiproxyProxies...)

	return proxies
}

func (s *Service) ValidateProxies(freeProxyList *[]Proxy) {
	wg := sync.WaitGroup{}

	for i := range *freeProxyList {
		wg.Add(1)
		proxy := &(*freeProxyList)[i]

		go func(i int, proxy *Proxy) {
			defer wg.Done()
			// proxy.Validate()
			proxy.Check()
		}(i, proxy)
	}

	wg.Wait()
}
