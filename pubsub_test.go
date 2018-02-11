package dockerinitiator

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	instance, err := PubSub()
	assert.NoError(t, err)
	defer func() { instance.Stop() }()

	assert.NoError(t, instance.Probe(10*time.Second))

	response, err := http.Get(fmt.Sprintf("http://%s/v1/projects/%s/topics", instance.GetHost(), instance.GetProject()))
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, 200, response.StatusCode)
}
