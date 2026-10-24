package server

import (
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/sway-agency/tv/bridge/internal/config"
)

func guard(settings config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		headers := writer.Header()
		headers.Set("Cache-Control", "no-store")
		headers.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		headers.Set("Strict-Transport-Security", "max-age=31536000")
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("Referrer-Policy", "no-referrer")
		if !hostAllowed(settings.Hosts, request.Host) {
			http.Error(writer, "unknown host", http.StatusMisdirectedRequest)
			return
		}
		if !originAllowed(settings.AllowedOrigins, request.Header.Get("Origin")) {
			http.Error(writer, "origin not allowed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func hostAllowed(hosts []string, hostHeader string) bool {
	host, _, err := net.SplitHostPort(hostHeader)
	if err != nil {
		host = hostHeader
	}
	host = strings.Trim(host, "[]")
	return slices.ContainsFunc(hosts, func(allowed string) bool {
		return strings.EqualFold(allowed, host)
	})
}

func originAllowed(allowed []string, origin string) bool {
	return origin == "" || slices.Contains(allowed, origin)
}

func lockoutKey(address string) string {
	ip := net.ParseIP(address)
	if ip == nil || ip.To4() != nil {
		return address
	}
	return ip.Mask(net.CIDRMask(64, 128)).String() + "/64"
}

func clientAddress(request *http.Request) string {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return request.RemoteAddr
	}
	return host
}
