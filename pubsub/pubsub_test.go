// +build integration

package dockerinitiator

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	docker "github.com/docker/docker/client"
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
		Image:       "google/cloud-sdk:latest",
		Cmd:         []string{"gcloud", "beta", "emulators", "pubsub", "start", "--host-port", "0.0.0.0:8262"},
		ExposedPort: "8262",
	})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	assert.NoError(t, err)
	inspectResp, _, err := client.ImageInspectWithRaw(context.Background(), "google/cloud-sdk:latest")
	assert.NoError(t, err)

	assert.Equal(t, inspectResp.ID, instance.Container().Image)
}

func TestPubSubCustomPort(t *testing.T) {
	instance, err := PubSub(PubSubConfig{
		Cmd:         []string{"--host=0.0.0.0", "--port=8263"},
		ExposedPort: "8263",
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
