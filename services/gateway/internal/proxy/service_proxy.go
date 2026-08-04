package proxy

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

// ServiceRegistry maps service names to their HTTP addresses.
type ServiceRegistry struct {
	services map[string]string // name → "host:port"
}

// NewServiceRegistry creates a registry from config.
// Example: {"payment": "localhost:9092", "merchant": "merchant:9092"}
func NewServiceRegistry(services map[string]string) *ServiceRegistry {
	return &ServiceRegistry{services: services}
}

// GetAddress returns the HTTP address of the named service.
func (r *ServiceRegistry) GetAddress(serviceName string) (string, error) {
	addr, ok := r.services[serviceName]
	if !ok {
		return "", fmt.Errorf("unknown service: %s", serviceName)
	}
	return addr, nil
}

// ProxyToService creates a Fiber handler that proxies the request to a downstream HTTP service.
func ProxyToService(registry *ServiceRegistry, serviceName string) fiber.Handler {
	return func(c fiber.Ctx) error {
		addr, err := registry.GetAddress(serviceName)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": fmt.Sprintf("%s is not available yet", serviceName),
			})
		}

		// Adjust docker hostname resolution if we are local but config says localhost
		// In a real prod setup, addr should be exactly what's in DNS
		if strings.HasPrefix(addr, "localhost:") && serviceName == "merchant" {
			// fallback in case they are running inside docker
			addr = "merchant:9092"
		}

		// Strip /api/v1 from the URL because downstream services are mounted at root
		targetPath := strings.TrimPrefix(c.OriginalURL(), "/api/v1")
		url := "http://" + addr + targetPath
		
		if err := proxy.Do(c, url); err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to proxy to %s: %v", serviceName, err),
			})
		}
		
		return nil
	}
}

// ExtractServiceName parses the route path "/api/v1/payment/..." and returns "payment".
func ExtractServiceName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}