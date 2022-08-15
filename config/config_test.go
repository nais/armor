package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_IsProtected(t *testing.T) {
	for _, test := range []struct {
		name     string
		rule     string
		expected bool
	}{
		{
			name:     "Changes to a non-protected rule is granted",
			rule:     "10",
			expected: false,
		},
		{
			name:     "Changes to a protected rule is denied",
			rule:     "1000",
			expected: true,
		},
		{
			name:     "Changes to a protected rule is denied",
			rule:     "2147483647",
			expected: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := os.Setenv("ARMOR_PROTECTED_RULES", "1000,2147483647")
			assert.NoError(t, err)
			cfg, err := SetupConfig()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, cfg.IsProtectedRule(test.rule))
		})
	}
}
