package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/makarivmri12/bug/pkg/models"
)

var (
	targetURL  = flag.String("target", "", "Target URL to scan")
	method     = flag.String("method", "GET", "HTTP method")
	headers    = flag.String("headers", "", "HTTP headers (comma-separated)")
	cookies    = flag.String("cookies", "", "HTTP cookies (comma-separated)")
	output     = flag.String("output", "json", "Output format (json, table, csv)")
	verbose    = flag.Bool("verbose", false, "Enable verbose output")
	version    = flag.Bool("version", false, "Show version")
)

const Version = "1.0.0"

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("HLFA CLI v%s\n", Version)
		os.Exit(0)
	}

	if *targetURL == "" {
		fmt.Fprintf(os.Stderr, "Error: --target is required\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("🔍 HLFA Scanner\n")
	fmt.Printf("Target: %s\n", *targetURL)
	fmt.Printf("Method: %s\n", *method)

	// Create target
	target := models.Target{
		URL:         *targetURL,
		Type:        models.TargetTypeWeb,
		WorkspaceID: "default",
		Headers:     parseHeaders(*headers),
		Cookies:     parseCookies(*cookies),
	}

	if *verbose {
		fmt.Printf("\n📋 Configuration:\n")
		fmt.Printf("  Headers: %v\n", target.Headers)
		fmt.Printf("  Cookies: %v\n", target.Cookies)
	}

	// Perform initial request
	fmt.Printf("\n📡 Sending request...\n")
	resp, err := performRequest(*method, target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Display results
	fmt.Printf("\n✅ Response received:\n")
	fmt.Printf("  Status: %d\n", resp.StatusCode)
	fmt.Printf("  Content-Length: %d\n", len(resp.Body))

	if *output == "json" {
		fmt.Printf("\n📊 Results (JSON):\n")
		fmt.Printf("%s\n", resp.Body[:min(500, len(resp.Body))])
	}
}

func parseHeaders(headerStr string) map[string]string {
	headers := make(map[string]string)
	if headerStr == "" {
		return headers
	}
	// Parse headers logic here
	return headers
}

func parseCookies(cookieStr string) map[string]string {
	cookies := make(map[string]string)
	if cookieStr == "" {
		return cookies
	}
	// Parse cookies logic here
	return cookies
}

func performRequest(method string, target models.Target) (*models.ResponseEvidence, error) {
	req, err := http.NewRequest(method, target.URL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range target.Headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := make([]byte, 0)
	headers := make(map[string]string)

	return &models.ResponseEvidence{
		Status:  resp.StatusCode,
		Headers: headers,
		Body:    string(body),
		TimeMS:  0,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
