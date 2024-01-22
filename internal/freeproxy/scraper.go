package freeproxy

import (
	"sync"

	"github.com/gocolly/colly"
)

type Scraper struct{}

func New() *Scraper {
	return &Scraper{}
}

func (s *Scraper) GetFreeProxyList() []Proxy {
	var freeProxyList []Proxy

	c := colly.NewCollector()

	c.OnHTML(".fpl-list tbody tr", func(e *colly.HTMLElement) {
		ip := e.ChildText("td:nth-child(1)")
		port := e.ChildText("td:nth-child(2)")
		code := e.ChildText("td:nth-child(3)")
		country := e.ChildText("td:nth-child(4)")
		anonimityLevel := e.ChildText("td:nth-child(5)")
		google := e.ChildText("td:nth-child(6)")
		https := e.ChildText("td:nth-child(7)")

		proxy := Proxy{
			IP:      ip,
			Port:    port,
			Code:    code,
			Country: country,
		}

		switch anonimityLevel {
		case "transparent":
			proxy.AnonymityLevel = Transparent
		case "anonymous":
			proxy.AnonymityLevel = Anonymousi
		case "elite proxy":
			proxy.AnonymityLevel = Elite
		default:
			proxy.AnonymityLevel = Transparent
		}

		if google == "yes" {
			proxy.Google = true
		} else {
			proxy.Google = false
		}

		if https == "yes" {
			proxy.HTTPS = true
		} else {
			proxy.HTTPS = false
		}

		freeProxyList = append(freeProxyList, proxy)
	})

	c.Visit("https://free-proxy-list.net/")

	return freeProxyList
}

func (s *Scraper) ValidateProxies(freeProxyList *[]Proxy) {
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
