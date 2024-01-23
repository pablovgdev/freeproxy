package main

import (
	"fmt"

	"github.com/pablovgdev/freeproxy/internal/models"
	"github.com/pablovgdev/freeproxy/internal/service"
)

func main() {
	s := service.Service{}
	params := service.GetProxiesParams{
		Uptime:         80,
		Protocol:       models.HTTPS,
		AnonimityLevel: models.Elite,
	}
	freeProxyList := s.GetProxies(params)
	s.ValidateProxies(&freeProxyList)

	for _, proxy := range freeProxyList {
		fmt.Println(proxy)
		// fmt.Println(scraper.ValidateProxy(proxy))
	}

}
