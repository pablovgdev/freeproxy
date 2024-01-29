package main

import (
	"fmt"
	"strconv"

	"github.com/pablovgdev/freeproxy/internal/freeproxy"
)

func main() {
	s := freeproxy.New()
	params := freeproxy.GetProxiesParams{
		Uptime:         80,
		Protocol:       freeproxy.HTTPS,
		AnonimityLevel: freeproxy.Elite,
	}
	freeProxyList := s.GetProxies(params)
	s.ValidateProxies(&freeProxyList)

	for _, proxy := range freeProxyList {
		if proxy.Valid {
			fmt.Println(proxy.IP + ":" + strconv.Itoa(proxy.Port))
		}
		// fmt.Println(scraper.ValidateProxy(proxy))
	}

}
