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

// MysqlConfig contains configs for mysql, User is automatically root
type MysqlConfig struct {
	Password     string
	DbName       string
	ProbeTimeout time.Duration
}

// Mysql starts up a mysql instance
func Mysql(config MysqlConfig) (*MysqlInstance, error) {

	if config.ProbeTimeout == 0 {
		config.ProbeTimeout = 10 * time.Second
	}

	i, err := dockerinitiator.CreateContainer(
		dockerinitiator.ContainerConfig{
			Image:         "storytel/mysql-57-test",
			Cmd:           []string{},
			Env:           []string{"MYSQL_ALLOW_EMPTY_PASSWORD=true", fmt.Sprintf("MYSQL_DATABASE=%s", config.DbName)},
			ContainerPort: "3306",
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
