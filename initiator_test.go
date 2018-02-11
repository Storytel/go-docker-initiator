package dockerinitiator

import (
	"log"
	"regexp"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := ClearObsolete(); err != nil {
		log.Panic(err)
	}

	m.Run()
}

func assertNumContainers(t *testing.T, num int) {
	assertNumContainersFilter(t, num, map[string][]string{})
}

func assertNumContainersFilter(t *testing.T, num int, filters map[string][]string) {
	client, err := docker.NewClientFromEnv()
	assert.NoError(t, err)

	containers, err := client.ListContainers(docker.ListContainersOptions{
		Filters: filters,
	})
	assert.NoError(t, err)

	assert.Len(t, containers, num)
}

func TestCreateContainer(t *testing.T) {
	instance, err := createContainer("ubuntu:latest", []string{"sleep", "300"}, "8080")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	assertNumContainersFilter(t, 1, map[string][]string{"id": []string{instance.container.ID}})
}

func TestTwoInstanceCoexist(t *testing.T) {
	instance1, err := createContainer("ubuntu:latest", []string{"sleep", "300"}, "8080")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, instance1.Stop())
	}()

	instance2, err := createContainer("ubuntu:latest", []string{"sleep", "300"}, "8080")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, instance2.Stop())
	}()

	assertNumContainers(t, 2)
}

func TestGetHost(t *testing.T) {
	instance, err := createContainer("ubuntu:latest", []string{"sleep", "300"}, "8080")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	assert.Regexp(t, regexp.MustCompile("^localhost:\\d+$"), instance.GetHost())
}

func TestClearObsolete(t *testing.T) {
	instance, err := createContainer("ubuntu:latest", []string{"sleep", "300"}, "8080")
	assert.NoError(t, err)
	defer func() { instance.Stop() }()

	OBSOLETE_AFTER = -9999 // So they're determined obsolete
	err = ClearObsolete()
	assert.NoError(t, err)

	assertNumContainers(t, 0)
}
