package dockerinitiator

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// ContainerConfig holds a subset of all of the configuration options available for the docker instance
type ContainerConfig struct {
	Image         string
	Cmd           []string
	Env           []string
	ContainerPort string
	Tmpfs         map[string]string
}

var ObsoleteAfter = 10 * time.Minute
var creator = "go-docker-initiator"

// CreateContainer applies the config and creates the container
func CreateContainer(config ContainerConfig, prober Probe) (*Instance, error) {
	ctx := context.Background()

	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	rc, err := client.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	io.Copy(ioutil.Discard, rc)
	defer rc.Close()

	createResp, err := client.ContainerCreate(ctx, &container.Config{
		Image: config.Image,
		Cmd:   config.Cmd,
		Env:   config.Env,
		ExposedPorts: nat.PortSet{
			nat.Port(config.ContainerPort): struct{}{},
		},

		Labels: map[string]string{"creator": creator},
	}, &container.HostConfig{
		PublishAllPorts: true,
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err = client.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	inspectResp, err := client.ContainerInspect(ctx, createResp.ID)
	if err != nil {
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", "127.0.0.1", getHostPort(inspectResp, config.ContainerPort))

	instance := &Instance{
		client:    client,
		host:      host,
		probe:     prober,
		container: inspectResp,
	}

	return instance, nil
}

func getHostPort(container types.ContainerJSON, containerport string) uint16 {
	if !strings.Contains(containerport, "/") {
		containerport += "/tcp" // It's the default
	}

	val, ok := container.NetworkSettings.Ports[nat.Port(containerport)]
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
	ctx := context.Background()

	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	apicontainers, err := client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "label", Value: fmt.Sprintf("creator=%s", creator)}),
		All:     true,
	})
	if err != nil {
		return err
	}

	for _, apicontainer := range apicontainers {
		inspectResp, err := client.ContainerInspect(ctx, apicontainer.ID)
		if err != nil {
			return err
		}

		startedAt, err := time.Parse(time.RFC3339, inspectResp.State.StartedAt)
		if err != nil {
			return err
		}

		if time.Since(startedAt) > ObsoleteAfter {
			log.Printf("Removing obsolete container %s", inspectResp.ID)
			err = client.ContainerRemove(ctx, inspectResp.ID, types.ContainerRemoveOptions{
				Force: true,
			})
			if err != nil {
				return err
			}
		}

	}

	return nil
}
