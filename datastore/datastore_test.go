package datastore_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	. "github.com/Storytel/go-docker-initiator/datastore"
	docker "github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestDatastore(t *testing.T) {
	instance, err := Datastore(DatastoreConfig{})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	response, err := http.Get(fmt.Sprintf("http://%s", instance.GetHost()))
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, 200, response.StatusCode)
}

func TestDatastoreCustomImage(t *testing.T) {
	instance, err := Datastore(DatastoreConfig{
		Image: "google/cloud-sdk:322.0.0-emulators",
		Cmd: []string{
			"/google-cloud-sdk/platform/cloud-datastore-emulator/cloud_datastore_emulator",
			"start",
			"--host=0.0.0.0",
			"--port=8263",
			"--store_on_disk=false",
		},
		ExposedPort: "8263",
	})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	assert.NoError(t, err)
	inspectResp, _, err := client.ImageInspectWithRaw(context.Background(), "google/cloud-sdk:322.0.0-emulators")
	assert.NoError(t, err)

	assert.Equal(t, inspectResp.ID, instance.Container().Image)
}

func TestDatastoreCustomPort(t *testing.T) {
	instance, err := Datastore(DatastoreConfig{
		Cmd:         []string{"--host=0.0.0.0", "--port=7263"},
		ExposedPort: "7263",
	})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	response, err := http.Get(fmt.Sprintf("http://%s", instance.GetHost()))
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, 200, response.StatusCode)
}
