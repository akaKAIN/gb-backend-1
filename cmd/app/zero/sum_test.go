package zero

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSum(t *testing.T) {
	result := Sum(1, 2)
	assert.Equal(t, 3, result, "Should be equal")
}
