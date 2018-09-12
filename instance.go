package dockerinitiator

import (
	"context"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

// Instance contains the setup for the docker instance
type Instance struct {
	client    *docker.Client
	host      string
	probe     Probe
	container *docker.Container
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

func (i *Instance) Container() *docker.Container {
	return i.container
}
