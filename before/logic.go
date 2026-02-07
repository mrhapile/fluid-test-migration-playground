package before

import (
	"fmt"
	"strings"
)

// Config represents a simplified configuration structure
// often seen in Fluid controllers (e.g. Dataset/Runtime specs).
type Config struct {
	Name      string
	Namespace string
	MountPath string
	Replicas  int
	Options   map[string]string
}

// ValidateAndNormalize checks if the configuration is valid
// and normalizes paths or defaults.
// This simulates business logic found in pkg/ddc/...
func ValidateAndNormalize(c *Config) error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if c.Replicas < 0 {
		return fmt.Errorf("replicas must be non-negative")
	}

	// Default logic mixed with validation
	if c.MountPath == "" {
		c.MountPath = "/mnt/data"
	}

	// Complex branching logic
	if strings.HasPrefix(c.MountPath, "//") {
		c.MountPath = strings.Replace(c.MountPath, "//", "/", 1)
	}

	// Check for restricted options
	if val, ok := c.Options["storage"]; ok {
		if val == "temporary" && c.Replicas > 1 {
			return fmt.Errorf("temporary storage does not support replicas > 1")
		}
	}

	return nil
}

// BuildConnectionURL simulates constructing a connection string
func BuildConnectionURL(c *Config) string {
	schema := "tcp"
	if val, ok := c.Options["ssl"]; ok && val == "true" {
		schema = "ssl"
	}
	return fmt.Sprintf("%s://%s.%s:8080", schema, c.Name, c.Namespace)
}
