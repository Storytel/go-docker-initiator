package dockerinitiator

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// MysqlInstance contains the config for mysql instance
type MysqlInstance struct {
	*Instance
	project string
}

// Mysql starts up a mysql instance
func Mysql(dbName string) (*MysqlInstance, error) {
	i, err := createContainer(
		ContainerConfig{
			Image:         "storytel/mysql-57-test",
			Cmd:           []string{},
			Env:           []string{"MYSQL_ALLOW_EMPTY_PASSWORD=true", fmt.Sprintf("MYSQL_DATABASE=%s", dbName)},
			ContainerPort: "3306",
			Tmpfs: map[string]string{
				"/var/lib/mysql": "rw",
			},
		},
		TCPProbe{})
	if err != nil {
		return nil, err
	}

	project := "__docker_initiator__project-" + strconv.Itoa(rand.Int())[:8]
	mi := &MysqlInstance{
		i,
		project,
	}

	if err = mi.Probe(10 * time.Second); err != nil {
		return nil, err
	}

	return mi, nil
}

// Setenv sets the required variables for running against the emulator
func (mi *MysqlInstance) Setenv(dbName string) error {
	err := os.Setenv("DB_SERVERNAME", mi.GetHost())
	if err != nil {
		return err
	}

	err = os.Setenv("DB_USERNAME", "root")
	if err != nil {
		return err
	}

	err = os.Setenv("DB_PASSWORD", "")
	if err != nil {
		return err
	}

	err = os.Setenv("DB_NAME", dbName)
	if err != nil {
		return err
	}

	return nil
}

// GetProject fetches the project for the mysql instance
func (mi *MysqlInstance) GetProject() string {
	return mi.project
}
