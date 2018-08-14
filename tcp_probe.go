package dockerinitiator

import (
	"net"
	"time"
)

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
