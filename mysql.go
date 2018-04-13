package dockerinitiator

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// MysqlInstance contains the config for mysql instance
type MysqlInstance struct {
	*Instance
	project string
	MysqlConfig
}

// MysqlConfig contains configs for mysql
type MysqlConfig struct {
	User     string
	Password string
	DbName   string
}

// Mysql starts up a mysql instance
func Mysql(config MysqlConfig) (*MysqlInstance, error) {
	i, err := createContainer(
		ContainerConfig{
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

	if err = mi.Probe(10 * time.Second); err != nil {
		return nil, err
	}

	return mi, nil
}

// Setenv sets the required variables for running against the emulator
func (mi *MysqlInstance) Setenv() error {
	if err := os.Setenv("DB_SERVERNAME", mi.GetHost()); err != nil {
		return err
	}

	if err := os.Setenv("DB_USERNAME", mi.User); err != nil {
		return err
	}

	if err := os.Setenv("DB_PASSWORD", mi.Password); err != nil {
		return err
	}

	return os.Setenv("DB_NAME", mi.DbName)
}

// SeedDatabase loads a seed sql file and executes it against the db.
func (mi *MysqlInstance) SeedDatabase(seedFilePath string) error {
	instanceHostAndPort := strings.Split(mi.GetHost(), ":")
	hostName := instanceHostAndPort[0]
	if hostName == "localhost" {
		hostName = "127.0.0.1"
	}
	hostPort := instanceHostAndPort[1]

	cmd := exec.Command("mysql", fmt.Sprintf("-h%s", hostName), fmt.Sprintf("-u%s", mi.User), fmt.Sprintf("-P%s", hostPort), mi.DbName, "-e", fmt.Sprintf("source %s", seedFilePath))

	var out, stderr bytes.Buffer

	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}

	return nil
}

// GetProject fetches the project for the mysql instance
func (mi *MysqlInstance) GetProject() string {
	return mi.project
}
