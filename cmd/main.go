package main

import (
	"fmt"
	"time"

	"github.com/pablovgdev/freeproxy/internal/freeproxy"
)

func main() {
	start := time.Now()
	s := freeproxy.New()
	params := freeproxy.GetProxiesParams{
		Uptime:         80,
		Protocol:       freeproxy.HTTPS,
		AnonimityLevel: freeproxy.Elite,
	}
	freeProxyList := s.GetProxies(params)
	s.ValidateProxies(&freeProxyList)

	for _, proxy := range freeProxyList {
		// proxy.Validate()
		if proxy.Validation.ValidTimes > 0 {
			fmt.Println(proxy.String())
			// fmt.Println(proxy.IP + ":" + strconv.Itoa(proxy.Port))
		}
		// fmt.Println(scraper.ValidateProxy(proxy))
	}
	fmt.Println("Total execution time: ", time.Since(start))
}
