package mysql

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	dockerinitiator "github.com/Storytel/go-docker-initiator"
)

// MysqlInstance contains the config for mysql instance
type MysqlInstance struct {
	*dockerinitiator.Instance
	project string
	MysqlConfig
}

var (
	DefaultImage = "storytel/mysql-57-test"

	defaultCmd = []string{}

	defaultExposedPort = "3306"
)

// MysqlConfig contains configs for mysql, User is automatically root
type MysqlConfig struct {

	// Password is the mysql password for the standard user "root"
	Password string

	// DbName is the name of the database you want to create and connect to
	DbName string

	// ProbeTimeout specifies the timeout for the probing.
	// A timeout results in a startup error, if left empty a default value is used
	ProbeTimeout time.Duration

	// Image specifies the image used for the Mysql docker instance.
	// If left empty it will be set to DefaultImage
	Image string
}

// Mysql starts up a mysql instance
func Mysql(config MysqlConfig) (*MysqlInstance, error) {

	if config.ProbeTimeout == 0 {
		config.ProbeTimeout = 10 * time.Second
	}

	if config.Image == "" {
		config.Image = DefaultImage
	}

	i, err := dockerinitiator.CreateContainer(
		dockerinitiator.ContainerConfig{
			Image:         config.Image,
			Cmd:           defaultCmd,
			Env:           []string{"MYSQL_ALLOW_EMPTY_PASSWORD=true", fmt.Sprintf("MYSQL_DATABASE=%s", config.DbName)},
			ContainerPort: defaultExposedPort,
			Tmpfs: map[string]string{
				"/var/lib/mysql": "rw",
			},
		},
		MysqlProbe{
			config,
		})
	if err != nil {
		return nil, err
	}

	project := "__docker_initiator__project-" + strconv.Itoa(rand.Int())[:8]
	mi := &MysqlInstance{
		i,
		project,
		config,
	}

	if err = mi.Probe(mi.ProbeTimeout); err != nil {
		return nil, err
	}

	return mi, nil
}

// Setenv sets the required variables for running against the emulator
func (mi *MysqlInstance) Setenv() error {
	if err := os.Setenv("MYSQL_SERVER", mi.GetHost()); err != nil {
		return err
	}

	if err := os.Setenv("MYSQL_USER", "root"); err != nil {
		return err
	}

	if err := os.Setenv("MYSQL_PASSWORD", mi.Password); err != nil {
		return err
	}

	return os.Setenv("MYSQL_DATABASE", mi.DbName)
}

// GetProject fetches the project for the mysql instance
func (mi *MysqlInstance) GetProject() string {
	return mi.project
}
