package dockerinitiator

import (
	"context"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

// Instance contains the setup for the docker instance
type Instance struct {
	client    *docker.Client
	container *docker.Container
	host      string
	probe     Probe
}

// Stop will remove the instance container
func (i *Instance) Stop() error {
	return i.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    i.container.ID,
		Force: true,
	})
}

// GetHost will fetch the host of the instance container
func (i *Instance) GetHost() string {
	return i.host
}

// Probe will periodically, during a timeout, check for an active connection in the instance container
func (i *Instance) Probe(timeout time.Duration) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	if err := i.probe.DoProbe(i); err == nil {
		return nil
	}

	for {
		select {
		case <-time.After(1 * time.Second):
			if err := i.probe.DoProbe(i); err == nil {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// GetDockerContainer returns the underlaying *docker.Container for advanced interaction
func (i *Instance) GetDockerContainer() *docker.Container {
	return i.container
}

// GetDockerClient returns the docker client used to control this instance
func (i *Instance) GetDockerClient() *docker.Client {
	return i.client
}
