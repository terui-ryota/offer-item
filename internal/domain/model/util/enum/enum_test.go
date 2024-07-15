package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	il := []int{1, 2, 3, 4}
	for _, i := range il {
		assert.Equal(t, true, Validate(i, il))
	}
	assert.Equal(t, false, Validate(5, il))
}
