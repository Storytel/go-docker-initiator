// +build integration

package dockerinitiator

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	instance, err := PubSub(PubSubConfig{})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	response, err := http.Get(fmt.Sprintf("http://%s/v1/projects/%s/topics", instance.GetHost(), instance.GetProject()))
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, 200, response.StatusCode)
}
