// +build integration

package dockerinitiator

import (
	"fmt"
	"net/http"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
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

	myClient, err := docker.NewClientFromEnv()
	assert.NoError(t, err)
	image, err := myClient.InspectImage("google/cloud-sdk:latest")
	assert.NoError(t, err)

	assert.Equal(t, image.ID, instance.Container().Image)
}

func TestPubSubCustomPort(t *testing.T) {
	instance, err := PubSub(PubSubConfig{
		Cmd:  []string{"--host=0.0.0.0", "--port=8263"},
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
