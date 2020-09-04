// +build integration

package mysql_test

import (
	"context"
	"net"
	"testing"
	"time"

	. "github.com/Storytel/go-docker-initiator/mysql"
	docker "github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestMysql(t *testing.T) {
	instance, err := Mysql(MysqlConfig{
		Password: "",
		DbName:   "test-db",
	})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	_, err = net.DialTimeout("tcp", instance.GetHost(), 1*time.Second)
	if !assert.NoError(t, err) {
		return
	}
}

func TestMysqlCustomImage(t *testing.T) {
	instance, err := Mysql(MysqlConfig{
		Password: "",
		DbName:   "test-db",
		Image:    "mysql:5.7",
	})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, instance.Stop())
	}()

	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	assert.NoError(t, err)
	image, _, err := client.ImageInspectWithRaw(context.Background(), "mysql:5.7")
	assert.NoError(t, err)

	assert.Equal(t, image.ID, instance.Container().Image)
}
