package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

var authHTTPClient = &http.Client{Timeout: 10 * time.Second}
func withScheme(addr string) string {
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return addr
	}
	return "http://" + addr
}

func AuthProxy(authAddr, targetPath string) fiber.Handler {
	baseURL := withScheme(authAddr)
	return func(c fiber.Ctx) error {
		upstreamPath := targetPath
		for _, p := range c.Route().Params {
			upstreamPath = strings.ReplaceAll(upstreamPath, ":"+p, c.Params(p))
		}
		body := c.Body()
		
		
		req, err := http.NewRequestWithContext(c.Context(), c.Method(), baseURL+upstreamPath, bytes.NewReader(body))
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"success": false,
				"error":   "failed to create upstream request",
			})
		}

		req.Header.Set("Content-Type", "application/json")
		if v := c.Get("X-Device-Id"); v != "" {
			req.Header.Set("X-Device-Id", v)
		}
		if v := c.Get("Authorization"); v != "" {
			req.Header.Set("Authorization", v)
		}
		if v := c.Get("User-Agent"); v != "" {
			req.Header.Set("User-Agent", v)
		}
		if uid, ok := c.Locals("user_id").(string); ok && uid != "" {
			req.Header.Set("X-User-Id", uid)
		}

		resp, err := authHTTPClient.Do(req)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("auth-svc unreachable: %v", err),
			})
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"success": false,
				"error":   "failed to read upstream response",
			})
		}

		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		return c.Status(resp.StatusCode).Send(respBody)
	}
}

// AuthProxyGET returns a Fiber handler that proxies GET requests to auth-svc.
func AuthProxyGET(authAddr, targetPath string) fiber.Handler {
	baseURL := withScheme(authAddr)
	
	return func(c fiber.Ctx) error {
		req, err := http.NewRequestWithContext(c.Context(), http.MethodGet, baseURL+targetPath, nil)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"success": false,
				"error":   "failed to create upstream request",
			})
		}

		req.Header.Set("Content-Type", "application/json")
		if v := c.Get("Authorization"); v != "" {
			req.Header.Set("Authorization", v)
		}

		resp, err := authHTTPClient.Do(req)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("auth-svc unreachable: %v", err),
			})
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"success": false,
				"error":   "failed to read upstream response",
			})
		}
		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		return c.Status(resp.StatusCode).Send(respBody)
	}
}
