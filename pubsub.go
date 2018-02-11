package dockerinitiator

import (
	"math/rand"
	"os"
	"strconv"
)

type PubSubInstance struct {
	*Instance
	project string
}

func PubSub() (*PubSubInstance, error) {
	i, err := createContainer(
		"google/cloud-sdk:latest",
		[]string{"gcloud", "beta", "emulators", "pubsub", "start", "--host-port=0.0.0.0:8262"},
		"8262",
	)
	if err != nil {
		return nil, err
	}

	project := "__storytel_initiator__project-" + strconv.Itoa(rand.Int())
	psi := &PubSubInstance{i, project}

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
