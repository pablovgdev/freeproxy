package freeproxy

import (
	"fmt"
	"strings"

	"crypto/tls"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
)

var agent *gorequest.SuperAgent
var requestPool = sync.Pool{
	New: func() interface{} {
		return agent.Clone()
	},
}
var resultPool = sync.Pool{
	New: func() interface{} {
		return map[string]interface{}{}
	},
}

func init() {
	agent = gorequest.New()
	agent.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	agent.Timeout(15 * time.Second)
}

var CheckURL = "https://httpbin.org/get"

func CheckProxy(proxy string) bool {
	fmt.Println(proxy)
	if strings.TrimSpace(proxy) == "" {
		return false
	}
	// no protocol, add //
	if !strings.Contains(proxy, "//") {
		proxy = "//" + proxy
	}
	// get resources from pool and release after operations
	request := requestPool.Get().(*gorequest.SuperAgent)
	resp := resultPool.Get().(map[string]interface{})
	defer requestPool.Put(request)
	defer resultPool.Put(resp)
	// do the Request
	_, _, errors := request.Proxy(proxy).Get(CheckURL).EndStruct(&resp)
	if errors != nil {
		return false
	}
	return strings.Contains(proxy, resp["origin"].(string))
}
