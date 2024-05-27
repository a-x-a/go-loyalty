package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderStatus_String(t *testing.T) {
	assert.Equal(t, "NEW", NEW.String())
	assert.Equal(t, "PROCESSING", PROCESSING.String())
	assert.Equal(t, "INVALID", INVALID.String())
	assert.Equal(t, "PROCESSED", PROCESSED.String())
}

func TestOrderStatus_Index(t *testing.T) {
	assert.Equal(t, 1, NEW.Index())
	assert.Equal(t, 2, PROCESSING.Index())
	assert.Equal(t, 3, INVALID.Index())
	assert.Equal(t, 4, PROCESSED.Index())
}
