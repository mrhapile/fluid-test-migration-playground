package before

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConnectionURL(t *testing.T) {
	testCases := []struct {
		desc   string
		config Config
		expect string
	}{
		{
			desc: "simple tcp url",
			config: Config{
				Name:      "test-dataset",
				Namespace: "default",
				Options:   map[string]string{},
			},
			expect: "tcp://test-dataset.default:8080",
		},
		{
			desc: "ssl enabled",
			config: Config{
				Name:      "secure-dataset",
				Namespace: "prod",
				Options:   map[string]string{"ssl": "true"},
			},
			expect: "ssl://secure-dataset.prod:8080",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := BuildConnectionURL(&tc.config)
			assert.Equal(t, tc.expect, result)
		})
	}
}

func TestValidateAndNormalize(t *testing.T) {
	// Table Driven approach with testify/assert
	tests := []struct {
		name        string
		cfg         *Config
		expectError bool
		expectedVal string
	}{
		{
			name:        "valid config",
			cfg:         &Config{Name: "valid", Replicas: 1},
			expectError: false,
			expectedVal: "/mnt/data",
		},
		{
			name:        "invalid empty name",
			cfg:         &Config{Replicas: 1},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAndNormalize(tt.cfg)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVal, tt.cfg.MountPath)
			}
		})
	}
}
