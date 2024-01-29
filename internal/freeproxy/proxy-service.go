package freeproxy

type GetProxiesParams struct {
	Language       string              `json:"language"`
	PageSize       int                 `json:"page_size"`
	Page           int                 `json:"page"`
	CountryCode    string              `json:"country_code"`
	AnonimityLevel ProxyAnonimityLevel `json:"anonymity"`
	Protocol       ProxyProtocol       `json:"protocol"`
	Speed          ProxySpeed          `json:"speed"`
	Uptime         int                 `json:"uptime"`
	GooglePassed   bool                `json:"google_passed"`
}

type ProxyService interface {
	GetProxies(params GetProxiesParams) []Proxy
}
