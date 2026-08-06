package netutil

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	// safeHTTPTimeout is the total deadline for a single outbound HTTP fetch.
	safeHTTPTimeout = 30 * time.Second
	// safeHTTPMaxBodySize caps the response body to prevent memory exhaustion
	// when a remote server streams an unexpectedly large payload.
	safeHTTPMaxBodySize = 10 << 20 // 10 MiB
)

// isBlockedIP reports whether ip points at an internal or reserved range
// that must not be reachable from a server-side fetch (SSRF mitigation).
func isBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	switch {
	case ip.IsLoopback(),
		ip.IsUnspecified(),
		ip.IsLinkLocalUnicast(),
		ip.IsLinkLocalMulticast(),
		ip.IsMulticast(),
		ip.IsPrivate(): // RFC 1918 (10/8, 172.16/12, 192.168/16) + RFC 4193 (fc00::/7)
		return true
	}
	return false
}

// safeTransport returns an http.Transport whose DialContext resolves the
// hostname and rejects the connection if any resolved IP is internal or
// reserved.  This closes the DNS-rebinding window at dial time.
func safeTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}
			for _, ipAddr := range ips {
				if isBlockedIP(ipAddr.IP) {
					return nil, fmt.Errorf("connection to internal/reserved IP blocked: %s", ipAddr.IP)
				}
			}
			return dialer.DialContext(ctx, network, addr)
		},
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}
}

// SafeHTTPClient returns an *http.Client configured with a bounded timeout
// and a transport that blocks connections to internal/reserved IP addresses.
// Callers should still validate the URL scheme via ValidateURL before use.
func SafeHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   safeHTTPTimeout,
		Transport: safeTransport(),
	}
}

// ValidateURL parses rawURL and ensures it uses an http(s) scheme.  It returns
// the parsed URL on success so the caller avoids a second parse.
func ValidateURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("missing host")
	}
	return u, nil
}

// LimitReader wraps r with the maximum response body size to prevent memory
// exhaustion from unbounded remote payloads.
func LimitReader(r io.Reader) io.Reader {
	return io.LimitReader(r, safeHTTPMaxBodySize)
}
