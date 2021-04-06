package dockerinitiator

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

// Instance contains the setup for the docker instance
type Instance struct {
	client    *docker.Client
	host      string
	probe     Probe
	container types.ContainerJSON
}

// Stop will remove the instance container
func (i *Instance) Stop() error {
	ctx := context.Background()
	return i.client.ContainerRemove(ctx, i.container.ID, types.ContainerRemoveOptions{Force: true})
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

func (i *Instance) Container() types.ContainerJSON {
	return i.container
}
