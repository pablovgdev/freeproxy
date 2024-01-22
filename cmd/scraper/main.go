package main

import (
	"fmt"

	freeproxy "github.com/pablovgdev/freeproxy/internal/freeproxy"
)

func main() {
	fp := freeproxy.New()
	freeProxyList := fp.GetFreeProxyList()
	fp.ValidateProxies(&freeProxyList)

	for _, proxy := range freeProxyList {
		fmt.Println(proxy)
		// fmt.Println(scraper.ValidateProxy(proxy))
	}

}
