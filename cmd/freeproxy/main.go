package main

import (
	"fmt"

	"github.com/pablovgdev/freeproxy/internal/freeproxy"
)

func main() {
	s := freeproxy.Service{}
	params := freeproxy.GetProxiesParams{
		Uptime:         80,
		Protocol:       freeproxy.HTTPS,
		AnonimityLevel: freeproxy.Elite,
	}
	freeProxyList := s.GetProxies(params)
	s.ValidateProxies(&freeProxyList)

	for _, proxy := range freeProxyList {
		fmt.Println(proxy)
		// fmt.Println(scraper.ValidateProxy(proxy))
	}

}
