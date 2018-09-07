// +build integration

package mysql

import (
	"net"
	"testing"
	"time"

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
	mi, err := Mysql(MysqlConfig{
		Password: "",
		DbName:   "test-db",
		Image:    "mysql:5.7",
	})
	if !assert.NoError(t, err) {
		return
	}

	defer func() {
		assert.NoError(t, mi.Stop())
	}()

	image, err := mi.GetDockerClient().InspectImage("mysql:5.7")
	assert.Equal(t, image.ID, mi.GetDockerContainer().Image)
}
