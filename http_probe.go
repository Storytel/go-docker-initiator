package dockerinitiator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var _ Probe = HTTPProbe{}

// HTTPProbe implements the IProbe interface for HTTP connections
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
	req = req.WithContext(reqctx)
	result, err := client.Do(req)
	if err != nil {
		return err
	}

	if result.StatusCode >= 200 && result.StatusCode < 300 {
		return nil
	}

	return errors.New("Invalid status: " + result.Status)
}
