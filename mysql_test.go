// +build integration

package dockerinitiator

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMysql(t *testing.T) {
	instance, err := Mysql(MysqlConfig{
		User:     "root",
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
