package main

import (
	"testing"
	"github.com/docker/docker/api/types/network"
	"time"
	"github.com/docker/docker/api/types"
	"io"
	"github.com/docker/docker/client"
	"errors"
)

func Test_buildServiceStates_AddsFirstSevenCharsOfGitHashToContainerName(t *testing.T) {
	moldConfig := &MoldConfig{
		Services: []DockerRunConfig{
			{
				Name: "servicename",
			},
		},
		LastCommit: "hash123456789",
	}

	networkConfig := &network.NetworkingConfig{}

	result, err := buildServiceStates(moldConfig, networkConfig)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	expected := "servicename-hash1234"
	if result[0].Name != expected {
		t.Errorf("expected %v but got %v", expected, result[0].Name)
	}
}

type remove func(containerID string, force bool) error

type mockDocker struct {
	removeFunc remove
}

func (mockDocker) BuildImageOfContainer(containerID string, reference string) error {
	panic("implement me")
}

func (mockDocker) BuildImageAsync(ic *ImageConfig, logWriter io.Writer, prefix string, done chan bool) error {
	panic("implement me")
}

func (mockDocker) Client() *client.Client {
	panic("implement me")
}

func (mockDocker) CreateNetwork(name string) (string, error) {
	panic("implement me")
}

func (mockDocker) GetImageList() ([]types.ImageSummary, error) {
	panic("implement me")
}

func (mockDocker) ImageAvailableLocally(imageName string) bool {
	panic("implement me")
}

func (md mockDocker) RemoveContainer(containerID string, force bool) error {
	return md.removeFunc(containerID, force)
}

func (mockDocker) RemoveImage(imageID string, force bool, cleanUp bool) error {
	panic("implement me")
}

func (mockDocker) RemoveNetwork(networkID string) error {
	panic("implement me")
}

func (mockDocker) PushImage(imageRef string, authCfg *types.AuthConfig, logWriter io.Writer, prefix string) error {
	panic("implement me")
}

func (mockDocker) StartContainer(cc *ContainerConfig, wr *Log, prefix string) error {
	panic("implement me")
}

func (mockDocker) StopContainer(containerID string, timeout time.Duration) error {
	panic("implement me")
}

func (mockDocker) TailLogs(containerID string, wr io.Writer, prefix string) error {
	panic("implement me")
}

func Test_removeContainers_RemovesAllByID(t *testing.T) {
	idsRemoved := make(map[string]bool)

	docker := mockDocker{
		removeFunc: func(containerID string, force bool) error {
			idsRemoved[containerID] = true
			return nil
		},
	}

	states := containerStates{
		&containerState{
			ContainerConfig: &ContainerConfig{
				id: "foo",
			},
		},
		&containerState{
			ContainerConfig: &ContainerConfig{
				id: "bar",
			},

		},
	}

	removeContainers(states, docker, nil)

	if !idsRemoved["foo"] {
		t.Error("Expected to have removed foo")
	}

	if !idsRemoved["bar"] {
		t.Error("Expected to have removed bar")
	}
}

func Test_removeContainers_AppendsErrors(t *testing.T) {
	docker := mockDocker{
		removeFunc: func(containerID string, force bool) error {
			return errors.New("bad")
		},
	}

	states := containerStates{
		&containerState{
			ContainerConfig: &ContainerConfig{
				id: "foo",
			},
		},
		&containerState{
			ContainerConfig: &ContainerConfig{
				id: "bar",
			},

		},
	}

	err := removeContainers(states, docker, nil)

	if err.Error() != "bad\nbad" {
		t.Errorf("Expected bad\\nbad but got %v", err)
	}
}