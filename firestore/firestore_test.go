package firestore_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	. "github.com/Storytel/go-docker-initiator/firestore"
	docker "github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestFirestore(t *testing.T) {
	instance, err := Firestore(FirestoreConfig{})
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

func TestFirestoreCustomImage(t *testing.T) {
	instance, err := Firestore(FirestoreConfig{
		Image:       "google/cloud-sdk:322.0.0-emulators",
		Cmd:         []string{"/google-cloud-sdk/platform/cloud-firestore-emulator/cloud_firestore_emulator", "--host=0.0.0.0", "--port=8263"},
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

func TestFirestoreCustomPort(t *testing.T) {
	instance, err := Firestore(FirestoreConfig{
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
