package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	require := require.New(t)

    l := InitLogger("info")
	require.NotNil(l)

	// defer func() {
    //     r := recover()
	// 	require.NotNil(r)
    // }()

    // InitLogger("invalid")
}
