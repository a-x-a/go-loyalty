package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderStatus_String(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("NEW", NEW.String())
	assert.Equal("PROCESSING", PROCESSING.String())
	assert.Equal("INVALID", INVALID.String())
	assert.Equal("PROCESSED", PROCESSED.String())
}

func TestOrderStatus_Index(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, NEW.Index())
	assert.Equal(2, PROCESSING.Index())
	assert.Equal(3, INVALID.Index())
	assert.Equal(4, PROCESSED.Index())
}
