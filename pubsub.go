package dockerinitiator

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

// PubSubInstance contains the instance config for a PubSub image
type PubSubInstance struct {
	*Instance
	project string
	PubSubConfig
}

// PubSubConfig contains configs for pubsub
type PubSubConfig struct {
	ProbeTimeout time.Duration
}

// PubSub will create a PubSub instance container
func PubSub(config PubSubConfig) (*PubSubInstance, error) {

	if config.ProbeTimeout == 0 {
		config.ProbeTimeout = 10 * time.Second
	}

	i, err := createContainer(
		ContainerConfig{
			Image:         "storytel/google-cloud-pubsub-emulator",
			Cmd:           []string{"--host=0.0.0.0", "-port=8262"},
			ContainerPort: "8262",
		},
		HTTPProbe{})
	if err != nil {
		return nil, err
	}

	project := "__docker_initiator__project-" + strconv.Itoa(rand.Int())[:8]
	psi := &PubSubInstance{
		i,
		project,
		config,
	}

	if err = psi.Probe(psi.ProbeTimeout); err != nil {
		return nil, err
	}

	return psi, nil
}

// Setenv sets the required variables for running against the emulator
func (psi *PubSubInstance) Setenv() error {
	err := os.Setenv("PUBSUB_EMULATOR_HOST", psi.GetHost())
	if err != nil {
		return err
	}

	err = os.Setenv("GOOGLE_CLOUD_PROJECT", psi.GetProject())
	if err != nil {
		return err
	}

	return nil
}

// GetProject fetches the project for the pubsub instance
func (psi *PubSubInstance) GetProject() string {
	return psi.project
}
