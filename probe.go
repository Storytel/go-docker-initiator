package dockerinitiator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//Probe provides an interfor for the probing mechanism
type Probe interface {
	DoProbe(instance *Instance) error
}

var _ Probe = TCPProbe{}

// TCPProbe implementes the IProbe interface for TPC connections
type TCPProbe struct {
}

// DoProbe will probe using TCP
func (i TCPProbe) DoProbe(instance *Instance) error {
	_, err := net.DialTimeout("tcp", instance.host, 1*time.Second)
	if err != nil {
		return err
	}

	return nil
}

var _ Probe = HTTPProbe{}

// HTTPProbe implementes the IProbe interface for HTTP connections
type HTTPProbe struct {
}

// DoProbe will probe using HTTP
func (i HTTPProbe) DoProbe(instance *Instance) error {
	url := fmt.Sprintf("http://%s/", instance.host)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	reqctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()
	req.WithContext(reqctx)
	result, err := client.Do(req)
	if err != nil {
		return err
	}

	if result.StatusCode >= 200 && result.StatusCode < 300 {
		return nil
	}

	return errors.New("Invalid status: " + result.Status)
}

// MysqlProbe implementes the IProbe interface for mysql instances
type MysqlProbe struct {
	MysqlConfig
}

// DoProbe will probe by waiting for log messages
func (i MysqlProbe) DoProbe(instance *Instance) error {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", i.User, i.Password, instance.GetHost(), i.DbName))
	defer db.Close()
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(0)

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}
