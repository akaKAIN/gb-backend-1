package zero

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	result := Sum(1, 2)
	assert.Equal(t, 3, result, "Should be equal")
}
