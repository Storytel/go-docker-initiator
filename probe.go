package dockerinitiator

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

//Probe provides an interfor for the probing mechanism
type Probe interface {
	DoProbe(host string) error
}

var _ Probe = TCPProbe{}

// TCPProbe implementes the IProbe interface for TPC connections
type TCPProbe struct {
}

// DoProbe will probe using TCP
func (i TCPProbe) DoProbe(host string) error {
	_, err := net.DialTimeout("tcp", host, 1*time.Second)
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
func (i HTTPProbe) DoProbe(host string) error {
	url := fmt.Sprintf("http://%s/", host)
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
