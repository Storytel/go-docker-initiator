package dockerinitiator_test

import (
	"context"
	"log"
	"os"
	"regexp"
	"testing"

	. "github.com/Storytel/go-docker-initiator"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ObsoleteAfter = -9999 // So they're determined obsolete
	if err := ClearObsolete(); err != nil {
		log.Panic(err)
	}

	m.Run()
	os.Exit(0)
}

func assertNumContainers(t *testing.T, num int) {
	assertNumContainersFilter(t, num, filters.NewArgs())
}

func assertNumContainersFilter(t *testing.T, num int, filters filters.Args) {
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	assert.NoError(t, err)

	filters.Add("label", "creator=go-docker-initiator")
	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters,
	})
	assert.NoError(t, err)

	assert.Len(t, containers, num)
}

func TestCreateContainer(t *testing.T) {
	instance, err := CreateContainer(
		ContainerConfig{
			Image:         "ubuntu:latest",
			Cmd:           []string{"sleep", "300"},
			ContainerPort: "8080",
		},
		HTTPProbe{})
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	assertNumContainersFilter(t, 1, filters.NewArgs(filters.Arg("id", instance.Container().ID)))
}

func TestTwoInstanceCoexist(t *testing.T) {
	instance1, err := CreateContainer(
		ContainerConfig{
			Image:         "ubuntu:latest",
			Cmd:           []string{"sleep", "300"},
			ContainerPort: "8080",
		},
		HTTPProbe{})
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, instance1.Stop())
	}()

	instance2, err := CreateContainer(
		ContainerConfig{
			Image:         "ubuntu:latest",
			Cmd:           []string{"sleep", "300"},
			ContainerPort: "8080",
		},
		HTTPProbe{})
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, instance2.Stop())
	}()

	assertNumContainers(t, 2)
}

func TestGetHost(t *testing.T) {
	instance, err := CreateContainer(
		ContainerConfig{
			Image:         "ubuntu:latest",
			Cmd:           []string{"sleep", "300"},
			ContainerPort: "8080",
		},
		HTTPProbe{})
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	assert.Regexp(t, regexp.MustCompile(`^127.0.0.1:\d+$`), instance.GetHost())
}

func TestClearObsolete(t *testing.T) {
	instance, err := CreateContainer(
		ContainerConfig{
			Image:         "ubuntu:latest",
			Cmd:           []string{"sleep", "300"},
			ContainerPort: "8080",
		},
		HTTPProbe{})
	assert.NoError(t, err)
	defer func() { instance.Stop() }()

	err = ClearObsolete()
	assert.NoError(t, err)

	assertNumContainers(t, 0)
}
