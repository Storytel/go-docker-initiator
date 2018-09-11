package dockerinitiator

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

// ContainerConfig holds a subset of all of the configuration options available for the docker instance
type ContainerConfig struct {
	Image         string
	Cmd           []string
	Env           []string
	ContainerPort string
	Tmpfs         map[string]string
}

var obsoleteAfter float64 = 10 * 60 // in seconds
var creator = "go-docker-initiator"

// CreateContainer applies the config and creates the container
func CreateContainer(config ContainerConfig, prober Probe) (*Instance, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	if err = client.PullImage(docker.PullImageOptions{Repository: config.Image}, docker.AuthConfiguration{}); err != nil {
		return nil, err
	}

	exposedports := map[docker.Port]struct{}{}
	exposedports[docker.Port(config.ContainerPort)] = struct{}{}
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		HostConfig: &docker.HostConfig{
			PublishAllPorts: true,
			Tmpfs:           config.Tmpfs,
		},
		Config: &docker.Config{
			Labels:       map[string]string{"creator": creator},
			Cmd:          config.Cmd,
			Env:          config.Env,
			ExposedPorts: exposedports,
			Image:        config.Image,
		},
	})
	if err != nil {
		return nil, err
	}

	if err = client.StartContainer(container.ID, &docker.HostConfig{}); err != nil {
		return nil, err
	}

	container, err = client.InspectContainer(container.ID)
	if err != nil {
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", "127.0.0.1", getHostPort(container, config.ContainerPort))

	instance := &Instance{
		client:    client,
		host:      host,
		probe:     prober,
		container: container,
	}

	return instance, nil
}

func getHostPort(container *docker.Container, containerport string) uint16 {
	if !strings.Contains(containerport, "/") {
		containerport += "/tcp" // It's the default
	}

	val, ok := container.NetworkSettings.Ports[docker.Port(containerport)]
	if !ok {
		log.Panic("No port configuration found on the created container")
	}

	port, err := strconv.ParseUint(val[0].HostPort, 10, 32)
	if err != nil {
		log.Panicf("Failed to parse the hostport (%s) to uint16", val[0].HostPort)
	}

	return uint16(port)
}

// ClearObsolete will delete all obsolete containers, created by this program
func ClearObsolete() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	apicontainers, err := client.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": []string{"creator=" + creator},
		},
		All: true,
	})
	if err != nil {
		return err
	}

	for _, apicontainer := range apicontainers {
		container, err := client.InspectContainer(apicontainer.ID)
		if err != nil {
			return err
		}
		startedAt := container.State.StartedAt
		if time.Since(startedAt).Seconds() > obsoleteAfter {
			log.Printf("Removing obsolete container %s", container.Name)
			err = client.RemoveContainer(docker.RemoveContainerOptions{
				ID:    container.ID,
				Force: true,
			})

			if err != nil {
				return err
			}
		}

	}

	return nil
}
