package ambassador

import (
	"errors"
	"io"

	"github.com/samalba/dockerclient"
)

var errNop = errors.New("Operation not supported")

// Ambassador is a wrapper around the Docker client that
// provides a shared volume and network for all containers.
type Ambassador struct {
	name   string
	client dockerclient.Client
}

// Create creates a new ambassador container.
func Create(client dockerclient.Client) (_ *Ambassador, err error) {
	amb := &Ambassador{
		client: client,
	}

	conf := &dockerclient.ContainerConfig{}
	host := &dockerclient.HostConfig{}
	conf.Entrypoint = []string{"/bin/sleep"}
	conf.Cmd = []string{"1d"}
	conf.Image = "busybox"
	conf.Volumes = map[string]struct{}{}
	conf.Volumes["/drone"] = struct{}{}

	// creates the ambassador container
	amb.name, err = client.CreateContainer(conf, "")
	if err != nil {

		// on failure attempts to pull the image
		client.PullImage(conf.Image, nil)

		// then attempts to re-create the container
		amb.name, err = client.CreateContainer(conf, "")
		if err != nil {
			return nil, err
		}
	}
	err = client.StartContainer(amb.name, host)
	return amb, err
}

// Destroy stops and deletes the ambassador container.
func (c *Ambassador) Destroy() error {
	c.client.StopContainer(c.name, 5)
	c.client.KillContainer(c.name, "SIGKILL")
	return c.client.RemoveContainer(c.name, true, true)
}

// CreateContainer creates a container.
func (c *Ambassador) CreateContainer(config *dockerclient.ContainerConfig, name string) (string, error) {
	return c.client.CreateContainer(config, name)
}

// InspectContainer returns container details.
func (c *Ambassador) InspectContainer(id string) (*dockerclient.ContainerInfo, error) {
	return c.client.InspectContainer(id)
}

// ContainerLogs returns an io.ReadCloser for reading the
// container logs.
func (c *Ambassador) ContainerLogs(id string, options *dockerclient.LogOptions) (io.ReadCloser, error) {
	return c.client.ContainerLogs(id, options)
}

// StartContainer starts a container. The ambassador volume
// is automatically linked. The ambassador network is linked
// iff a network mode is not already specified.
func (c *Ambassador) StartContainer(id string, config *dockerclient.HostConfig) error {
	config.VolumesFrom = append(config.VolumesFrom, "container:"+c.name)
	if len(config.NetworkMode) == 0 {
		config.NetworkMode = "container:" + c.name
	}
	return c.client.StartContainer(id, config)
}

// StopContainer stops a container.
func (c *Ambassador) StopContainer(id string, timeout int) error {
	return c.client.StopContainer(id, timeout)
}

// PullImage pulls an image.
func (c *Ambassador) PullImage(name string, auth *dockerclient.AuthConfig) error {
	return c.client.PullImage(name, auth)
}

// RemoveContainer removes a container.
func (c *Ambassador) RemoveContainer(id string, force, volumes bool) error {
	return c.client.RemoveContainer(id, force, volumes)
}

// KillContainer kills a running container.
func (c *Ambassador) KillContainer(id, signal string) error {
	return c.client.KillContainer(id, signal)
}

//
// methods below are not implemented
//

// Info returns a no-op error
func (c *Ambassador) Info() (*dockerclient.Info, error) {
	return nil, errNop
}

// ListContainers returns a no-op error
func (c *Ambassador) ListContainers(all bool, size bool, filters string) ([]dockerclient.Container, error) {
	return nil, errNop
}

// RestartContainer returns a no-op error
func (c *Ambassador) RestartContainer(id string, timeout int) error {
	return errNop
}

// StartMonitorEvents returns a no-op error
func (c *Ambassador) StartMonitorEvents(cb dockerclient.Callback, ec chan error, args ...interface{}) {

}

// StopAllMonitorEvents returns a no-op error
func (c *Ambassador) StopAllMonitorEvents() {

}

// Version returns a no-op error
func (c *Ambassador) Version() (*dockerclient.Version, error) {
	return nil, errNop
}

// ListImages returns a no-op error
func (c *Ambassador) ListImages() ([]*dockerclient.Image, error) {
	return nil, errNop
}

// RemoveImage returns a no-op error
func (c *Ambassador) RemoveImage(name string) error {
	return errNop
}

// PauseContainer returns a no-op error
func (c *Ambassador) PauseContainer(name string) error {
	return errNop
}

// UnpauseContainer returns a no-op error
func (c *Ambassador) UnpauseContainer(name string) error {
	return errNop
}

// Exec returns a no-op error
func (c *Ambassador) Exec(config *dockerclient.ExecConfig) (string, error) {
	var empty string
	return empty, errNop
}
