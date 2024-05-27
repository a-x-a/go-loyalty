package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Check(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input string
		want  bool
	}{
		{"5062821234567892", true},
		{"5062821734567892", false},
		{"5062 8212 3456 7892", false},
		{"123456789012345678901234567891", true},
		{"123456789012345678901234567890", false},
	}

	for _, tt := range tests {
		got := Check(tt.input)
		assert.Equal(tt.want, got)
	}
}

func Test_revers(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input string
		want  string
	}{
		{"1234567890", "0987654321"},
		{"gopher", "rehpog"},
		{"", ""},
	}

	for _, tt := range tests {
		got := revers(tt.input)
		assert.Equal(tt.want, got)
	}
}
