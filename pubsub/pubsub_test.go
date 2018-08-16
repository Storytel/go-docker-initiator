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

func TestPubSubCustomImage(t *testing.T) {
	instance, err := PubSub(PubSubConfig{
		Image: "google/cloud-sdk:latest",
		Cmd:   []string{"gcloud", "beta", "emulators", "pubsub", "start", "--host-port", "0.0.0.0:8262"},
		Port:  "8262",
	})
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

func TestPubSubCustomPort(t *testing.T) {
	instance, err := PubSub(PubSubConfig{
		Port: "8263",
	})
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
