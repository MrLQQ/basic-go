package demo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAbc(t *testing.T) {
	assert.True(t, "abc" > "bcd")
	t.Log("输出一句话")
	require.True(t, "abc" > "bcd")
	t.Log("再次输出一句话")
}
