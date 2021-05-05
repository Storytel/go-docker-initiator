package firestore

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	dockerinitiator "github.com/Storytel/go-docker-initiator"
)

// FirestoreInstance contains the instance config for a Firestore image
type FirestoreInstance struct {
	*dockerinitiator.Instance
	project string
	FirestoreConfig
}

var (
	DefaultImage = "storytel/gcp-firestore-emulator"

	DefaultCmd = []string{"--host=0.0.0.0", "--port=8263"}

	DefaultExposedPort = "8263"
)

// FirestoreConfig contains configs for firestore
type FirestoreConfig struct {
	// ProbeTimeout specifies the timeout for the probing.
	// A timeout results in a startup error, if left empty a default value is used
	ProbeTimeout time.Duration

	// Image specifies the image used for the Mysql docker instance.
	// If left empty it will be set to DefaultImage
	Image string

	// Cmd is the commands that will run in the container
	// Is left empty it will be set to DefaultCmd
	Cmd []string

	// ExposedPort sets the exposed port of the container
	// If left empty it will be set to DefaultExposedPort
	ExposedPort string
}

// Firestore will create a Firestore instance container
func Firestore(config FirestoreConfig) (*FirestoreInstance, error) {
	if config.ProbeTimeout == 0 {
		config.ProbeTimeout = 10 * time.Second
	}

	if config.Image == "" {
		config.Image = DefaultImage
	}

	if config.ExposedPort == "" {
		config.ExposedPort = DefaultExposedPort
	}

	if len(config.Cmd) == 0 {
		config.Cmd = DefaultCmd
	}

	i, err := dockerinitiator.CreateContainer(
		dockerinitiator.ContainerConfig{
			Image:         config.Image,
			Cmd:           config.Cmd,
			ContainerPort: config.ExposedPort,
		},
		dockerinitiator.HTTPProbe{})
	if err != nil {
		return nil, err
	}

	project := "__docker_initiator__project-" + strconv.Itoa(rand.Int())[:8]
	fsi := &FirestoreInstance{
		i,
		project,
		config,
	}

	if err = fsi.Probe(fsi.ProbeTimeout); err != nil {
		return nil, err
	}

	return fsi, nil
}

// Setenv sets the required variables for running against the emulator
func (fsi *FirestoreInstance) Setenv() error {
	err := os.Setenv("FIRESTORE_EMULATOR_HOST", fsi.GetHost())
	if err != nil {
		return err
	}

	err = os.Setenv("FIRESTORE_GOOGLE_CLOUD_PROJECT", fsi.GetProject())
	if err != nil {
		return err
	}

	return nil
}

// GetProject fetches the project for the firestore instance
func (fsi *FirestoreInstance) GetProject() string {
	return fsi.project
}
