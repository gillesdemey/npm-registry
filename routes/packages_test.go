package routes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryEscape(t *testing.T) {
	expected := "@foo%2Fbar"
	output := queryEscape("@foo/bar")
	assert.Equal(t, output, expected)
}
