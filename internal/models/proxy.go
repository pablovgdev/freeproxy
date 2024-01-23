package models

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Proxy struct {
	IP             string              `json:"ip"`
	Port           int                 `json:"port"`
	AnonymityLevel ProxyAnonimityLevel `json:"anonymity"`
	Protocol       ProxyProtocol       `json:"protocol"`
	Speed          ProxySpeed          `json:"speed"`
	Uptime         int                 `json:"uptime"`
	Latency        int                 `json:"latency"`
	GooglePassed   int                 `json:"google_passed"`
	CountryCode    string              `json:"country_code"`
	Valid          bool
	ResponseTime   int
}

type ProxyAnonimityLevel int

const (
	// Reveals IP and proxy usage.
	Transparent ProxyAnonimityLevel = 0
	// Hides IP but reveals proxy usage.
	Anonymous ProxyAnonimityLevel = 1
	// Hides IP and proxy usage.
	Elite ProxyAnonimityLevel = 2
)

type ProxyProtocol int

const (
	HTTP   ProxyProtocol = 1
	HTTPS  ProxyProtocol = 2
	Socks4 ProxyProtocol = 4
	Socks5 ProxyProtocol = 8
)

type ProxySpeed int

const (
	Slow   ProxySpeed = 0
	Medium ProxySpeed = 1
	Fast   ProxySpeed = 2
)

type ValidateProxyResponse struct {
	Origin string `json:"origin"`
}

func (p *Proxy) Validate() {
	address, err := url.Parse("http://" + p.IP + ":" + strconv.Itoa(p.Port))
	if err != nil {
		return
	}

	transport := &http.Transport{Proxy: http.ProxyURL(address)}
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	start := time.Now()
	resp, err := client.Get("https://httpbin.org/get")
	p.ResponseTime = int(time.Since(start).Seconds())

	if err != nil {
		log.Print(err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

	var validateProxyResponse ValidateProxyResponse
	err = json.Unmarshal(body, &validateProxyResponse)

	if err != nil {
		log.Fatal(err)
		return
	}

	respIP := validateProxyResponse.Origin

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
