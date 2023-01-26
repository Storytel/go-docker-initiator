package datastore

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	dockerinitiator "github.com/Storytel/go-docker-initiator"
)

// DatastoreInstance contains the instance config for a Datastore image
type DatastoreInstance struct {
	*dockerinitiator.Instance
	project string
	DatastoreConfig
}

var (
	DefaultImage = "storytel/google-cloud-datastore-emulator"

	DefaultCmd = []string{"--host=0.0.0.0", "--port=8263"}

	DefaultExposedPort = "8263"
)

// DatastoreConfig contains configs for datastore
type DatastoreConfig struct {
	// ProbeTimeout specifies the timeout for the probing.
	// A timeout results in a startup error, if left empty a default value is used
	ProbeTimeout time.Duration

	// Image specifies the image used for the Datastore docker instance.
	// If left empty it will be set to DefaultImage
	Image string

	// Cmd is the commands that will run in the container
	// Is left empty it will be set to DefaultCmd
	Cmd []string

	// ExposedPort sets the exposed port of the container
	// If left empty it will be set to DefaultExposedPort
	ExposedPort string
}

// Datastore will create a Datastore instance container
func Datastore(config DatastoreConfig) (*DatastoreInstance, error) {
	if config.ProbeTimeout == 0 {
		config.ProbeTimeout = 60 * time.Second
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

	project := "docker-initiator--project-" + strconv.Itoa(rand.Int())[:8]
	fsi := &DatastoreInstance{
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
func (fsi *DatastoreInstance) Setenv() error {
	err := os.Setenv("DATASTORE_EMULATOR_HOST", fsi.GetHost())
	if err != nil {
		return err
	}

	err = os.Setenv("GOOGLE_CLOUD_PROJECT", fsi.GetProject())
	if err != nil {
		return err
	}

	return nil
}

// GetProject fetches the project for the datastore instance
func (fsi *DatastoreInstance) GetProject() string {
	return fsi.project
}
