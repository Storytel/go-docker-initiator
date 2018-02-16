package dockerinitiator

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

type PubSubInstance struct {
	*Instance
	project string
}

func PubSub() (*PubSubInstance, error) {
	i, err := createContainer(
		"storytel/google-cloud-pubsub-emulator",
		[]string{"--host=0.0.0.0", "-port=8262"},
		"8262",
	)
	if err != nil {
		return nil, err
	}

	project := "__docker_initiator__project-" + strconv.Itoa(rand.Int())[:8]
	psi := &PubSubInstance{i, project}

	if err = psi.Probe(10 * time.Second); err != nil {
		return nil, err
	}

	return psi, nil
}

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

func (psi *PubSubInstance) GetProject() string {
	return psi.project
}
