package options_test

import (
	"github.com/shady831213/jarvism/core/options"
	"testing"
)

func TestOptionUsage(t *testing.T) {
	options.GetJvsOptions().Usage()
}
