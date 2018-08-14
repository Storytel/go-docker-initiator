package dockerinitiator

//Probe provides an interface for the probing mechanism
type Probe interface {
	DoProbe(instance *Instance) error
}
