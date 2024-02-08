package freeproxy

import (
	"sync"

	"github.com/pablovgdev/freeproxy/internal/providers"
	"github.com/pablovgdev/freeproxy/internal/proxy"
)

type FreeProxy struct {
	Providers []proxy.ProxyProvider
}

func NewFreeProxy() *FreeProxy {
	return &FreeProxy{
		Providers: []proxy.ProxyProvider{
			&providers.Lumiproxy{},
			&providers.ProxyScrape{},
		},
	}
}

func (fp *FreeProxy) GetProxies(params proxy.GetProxiesParams) []proxy.Proxy {
	proxies := []proxy.Proxy{}
	wg := sync.WaitGroup{}
	wg.Add(2)

	for _, provider := range fp.Providers {
		go func(provider proxy.ProxyProvider) {
			defer wg.Done()
			providerProxies := provider.GetProxies(params)
			proxies = append(proxies, providerProxies...)
		}(provider)
	}

	wg.Wait()

	return proxies
}

func (fp *FreeProxy) ValidateProxies(freeProxyList *[]proxy.Proxy) {
	wg := sync.WaitGroup{}

	for i := range *freeProxyList {
		wg.Add(1)
		item := &(*freeProxyList)[i]

		go func(i int, item *proxy.Proxy) {
			defer wg.Done()
			item.Validate()
		}(i, item)
	}

	wg.Wait()
}
