package freeproxy

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Proxy struct {
	IP             string
	Port           string
	Code           string
	Country        string
	AnonymityLevel ProxyAnonimityLevel
	Google         bool
	HTTPS          bool
	Valid          bool
	ResponseTime   uint8
}

type ProxyAnonimityLevel uint8

const (
	// Reveals IP and proxy usage.
	Transparent ProxyAnonimityLevel = 1
	// Hides IP but reveals proxy usage.
	Anonymousi ProxyAnonimityLevel = 2
	// Hides IP and proxy usage.
	Elite ProxyAnonimityLevel = 3
)

func (p *Proxy) Validate() {
	address, err := url.Parse("http://" + p.IP + ":" + p.Port)
	if err != nil {
		return
	}

	transport := &http.Transport{Proxy: http.ProxyURL(address)}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	client := &http.Client{Transport: transport, Timeout: 60 * time.Second}

	start := time.Now()
	resp, err := client.Get("http://ident.me")
	p.ResponseTime = uint8(time.Since(start).Seconds())

	if err != nil || resp.StatusCode != 200 {
		return
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	respIP := string(respBody)

	if len(respIP) == 0 || len(respIP) > 15 {
		return
	}

	partsIP := strings.Split(respIP, ".")

	if len(partsIP) < 3 {
		return
	}

	IP := partsIP[0] + "." + partsIP[1] + "." + partsIP[2]

	p.Valid = strings.Contains(p.IP, IP)
}
