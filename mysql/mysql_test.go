// +build integration

package mysql

import (
	"net"
	"testing"
	"time"

	"github.com/fsouza/go-dockerclient"

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

	myClient, err := docker.NewClientFromEnv()
	assert.NoError(t, err)
	image, err := myClient.InspectImage("mysql:5.7")
	assert.NoError(t, err)

	assert.Equal(t, image.ID, instance.Container().Image)
}

func TestMysqlCustomPort(t *testing.T) {
	instance, err := Mysql(MysqlConfig{
		Password: "",
		DbName:   "test-db",
		Image:    "mysql:5.7",
		Port:     "3306",
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
