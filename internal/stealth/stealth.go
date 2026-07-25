package stealth

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// UserAgents contains a list of common user agents
var UserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
}

// HeaderSpoofing handles HTTP header manipulation for stealth
type HeaderSpoofing struct {
	userAgents []string
}

// NewHeaderSpoofing creates a new header spoofing instance
func NewHeaderSpoofing() *HeaderSpoofing {
	return &HeaderSpoofing{
		userAgents: UserAgents,
	}
}

// ApplyStealth modifies request headers for stealth
func (hs *HeaderSpoofing) ApplyStealth(req *http.Request) {
	// Random user agent
	req.Header.Set("User-Agent", hs.getRandomUserAgent())

	// Add realistic referer
	req.Header.Set("Referer", "https://www.google.com/")

	// Add accept headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
}

// getRandomUserAgent returns a random user agent
func (hs *HeaderSpoofing) getRandomUserAgent() string {
	return hs.userAgents[rand.Intn(len(hs.userAgents))]
}

// ProxyRotator manages proxy rotation for distributed scanning
type ProxyRotator struct {
	proxies []string
	index   int
}

// NewProxyRotator creates a new proxy rotator
func NewProxyRotator(proxies []string) *ProxyRotator {
	rand.Seed(time.Now().UnixNano())
	return &ProxyRotator{
		proxies: proxies,
		index:   0,
	}
}

// GetNextProxy returns the next proxy in rotation
func (pr *ProxyRotator) GetNextProxy() string {
	if len(pr.proxies) == 0 {
		return ""
	}

	proxy := pr.proxies[pr.index]
	pr.index = (pr.index + 1) % len(pr.proxies)
	return proxy
}

// GetProxyURL formats a proxy URL
func GetProxyURL(proxyAddr string) string {
	return fmt.Sprintf("http://%s", proxyAddr)
}
