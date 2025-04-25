package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToInterfaceSlice(t *testing.T) {
	t.Run("converts string slice", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		result := ConvertToInterfaceSlice(input)
		
		assert.Equal(t, len(input), len(result), "output slice should have same length as input")
		for i, v := range input {
			assert.Equal(t, v, result[i], "values should match at index %d", i)
		}
	})

	t.Run("converts integer slice", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := ConvertToInterfaceSlice(input)
		
		assert.Equal(t, len(input), len(result), "output slice should have same length as input")
		for i, v := range input {
			assert.Equal(t, v, result[i], "values should match at index %d", i)
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		var input []string
		result := ConvertToInterfaceSlice(input)
		
		assert.Equal(t, 0, len(result), "output slice should be empty")
		assert.NotNil(t, result, "output slice should not be nil")
	})
} 