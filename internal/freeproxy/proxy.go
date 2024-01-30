package freeproxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Proxy struct {
	IP             string              `json:"ip"`
	Port           int                 `json:"port"`
	AnonymityLevel ProxyAnonimityLevel `json:"anonymity"`
	Protocol       ProxyProtocol       `json:"protocol"`
	CountryCode    string              `json:"country_code"`
	Provider       string              `json:"provider"`
	Validation     ProxyValidation     `json:"validation"`
}

type ProxyValidation struct {
	mu                 *sync.Mutex
	ValidationAttempts int `json:"validation_attempts"`
	ValidTimes         int `json:"valid_times"`
}

type ProxyAnonimityLevel int

const (
	// Reveals IP and proxy usage.
	Transparent ProxyAnonimityLevel = 1
	// Hides IP but reveals proxy usage.
	Anonymous ProxyAnonimityLevel = 2
	// Hides IP and proxy usage.
	Elite ProxyAnonimityLevel = 3
)

type ProxyProtocol int

const (
	HTTP   ProxyProtocol = 1
	HTTPS  ProxyProtocol = 2
	Socks4 ProxyProtocol = 3
	Socks5 ProxyProtocol = 4
)

type ProxySpeed int

const (
	Slow   ProxySpeed = 1
	Medium ProxySpeed = 2
	Fast   ProxySpeed = 3
)

type ValidateProxyResponse struct {
	Origin string `json:"origin"`
}

func NewProxy(
	IP string,
	Port int,
	AnonymityLevel ProxyAnonimityLevel,
	Protocol ProxyProtocol,
	CountryCode string,
	Provider string,
) *Proxy {
	return &Proxy{
		IP:             IP,
		Port:           Port,
		AnonymityLevel: AnonymityLevel,
		Protocol:       Protocol,
		CountryCode:    CountryCode,
		Provider:       Provider,
		Validation: ProxyValidation{
			mu:                 &sync.Mutex{},
			ValidationAttempts: 0,
			ValidTimes:         0,
		},
	}
}

func (p *Proxy) String() string {
	anonimity := ""

	switch p.AnonymityLevel {
	case Transparent:
		anonimity = "Transparent"
	case Anonymous:
		anonimity = "Anonymous"
	case Elite:
		anonimity = "Elite"
	default:
		anonimity = "Transparent"
	}

	protocol := ""

	switch p.Protocol {
	case HTTP:
		protocol = "HTTP"
	case HTTPS:
		protocol = "HTTPS"
	case Socks4:
		protocol = "Socks4"
	case Socks5:
		protocol = "Socks5"
	default:
		protocol = "HTTP"
	}

	uptime := p.Validation.ValidTimes * 100 / p.Validation.ValidationAttempts

	return fmt.Sprintf("%s:%d (%s %s %s %s %d%%)", p.IP, p.Port, p.Provider, anonimity, protocol, p.CountryCode, uptime)
}

func (p *Proxy) Validate() {
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			err := p.check()

			p.Validation.mu.Lock()
			defer p.Validation.mu.Unlock()
			p.Validation.ValidationAttempts++

			if err == nil {
				p.Validation.ValidTimes++
			}
		}()
	}

	wg.Wait()
}

func (p *Proxy) check() error {
	address, err := url.Parse("http://" + p.IP + ":" + strconv.Itoa(p.Port))
	if err != nil {
		return err
	}

	transport := &http.Transport{Proxy: http.ProxyURL(address)}
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	resp, err := client.Get("https://httpbin.org/ip")

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var validateProxyResponse ValidateProxyResponse
	err = json.Unmarshal(body, &validateProxyResponse)

	if err != nil {
		return err
	}

	respIP := validateProxyResponse.Origin

	if len(respIP) == 0 || len(respIP) > 15 {
		return errors.New("invalid IP format")
	}

	partsIP := strings.Split(respIP, ".")

	if len(partsIP) < 3 {
		return errors.New("invalid IP parts")
	}

	IP := partsIP[0] + "." + partsIP[1] + "." + partsIP[2]

	isValid := strings.Contains(p.IP, IP)

	if !isValid {
		return errors.New("IP does not match")
	}

	return nil
}
