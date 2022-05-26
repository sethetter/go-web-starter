package actions_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	s, _, _, _ := makeServer(t)
	defer s.Close()

	body, _ := sendRequest(t, s.URL, nil)
	assert.Contains(t, string(body), "Welcome")
}
