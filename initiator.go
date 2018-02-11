package dockerinitiator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

type Instance struct {
	client    *docker.Client
	container *docker.Container
	host      string
}

var OBSOLETE_AFTER float64 = 10 * 60 // in seconds
var CREATOR = "go-docker-initiator"

func createContainer(image string, cmd []string, containerport string) (*Instance, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	if err = client.PullImage(docker.PullImageOptions{Repository: image}, docker.AuthConfiguration{}); err != nil {
		return nil, err
	}

	exposedports := map[docker.Port]struct{}{}
	exposedports[docker.Port(containerport)] = struct{}{}
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		HostConfig: &docker.HostConfig{
			PublishAllPorts: true,
		},
		Config: &docker.Config{
			Labels:       map[string]string{"creator": CREATOR},
			Cmd:          cmd,
			ExposedPorts: exposedports,
			Image:        image,
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

	host := fmt.Sprintf("%s:%d", "localhost", getHostPort(container, containerport))

	instance := &Instance{client, container, host}

	return instance, nil
}

func (i *Instance) Stop() error {
	return i.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    i.container.ID,
		Force: true,
	})
}

func (i *Instance) GetHost() string {
	return i.host
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
		log.Panic("Failed to parse the hostport (%s) to uint16", val[0].HostPort)
	}

	return uint16(port)
}

func ClearObsolete() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	apicontainers, err := client.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": []string{"creator=" + CREATOR},
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
		if time.Since(startedAt).Seconds() > OBSOLETE_AFTER {
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

func (i *Instance) Probe(timeout time.Duration) error {
	url := fmt.Sprintf("http://%s/", i.GetHost())
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	doProbe := func() error {
		result, err := http.Get(url)
		if err != nil {
			return err
		}

		if result.StatusCode >= 200 && result.StatusCode < 300 {
			return nil
		}

		return errors.New("Invalid statuscode: " + strconv.Itoa(result.StatusCode))
	}

	if err := doProbe(); err == nil {
		return nil
	}

	for {
		select {
		case <-time.After(1 * time.Second):
			if err := doProbe(); err == nil {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
