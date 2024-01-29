package instance

import (
	"log"
	de "shedock/pkg/docker/client"
	"sync"
)

type instance struct {
	dockerClientInstance *de.Container
}

var singleton = &instance{}
var once sync.Once

func Init() {
	once.Do(func() {
		dockerClient, err := de.NewContainer()
		if err != nil {
			log.Fatalf("Failed to create container: %v", err)
		}
		singleton.dockerClientInstance = dockerClient
	})
}

func GetDockerInstance() *de.Container {
	return singleton.dockerClientInstance
}

// Destroy closes the connections & cleans up the instance
func Destroy() error {
	err := singleton.dockerClientInstance.DockerClient.Close()
	if err != nil {
		return err
	}
	return nil
}
