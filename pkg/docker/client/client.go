package client

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	imageName   = "alpine:3.18"
	alpineImage = "alpine"
)

type Container struct {
	// ContainerID is the ID of the container that will be created
	ContainerID string
	// DockerClient is the Docker client
	DockerClient *client.Client
	ctx          context.Context
}

func NewContainer() (*Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	err = pullImage(cli, imageName)
	if err != nil {
		return nil, fmt.Errorf("failed to pull Alpine image: %s", err)
	}

	container := &Container{
		DockerClient: cli,
		ctx:          context.Background(),
	}

	err = container.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %s", err)
	}

	return container, nil
}

func (c *Container) Create() error {
	resp, err := c.DockerClient.ContainerCreate(c.ctx, &container.Config{
		Image:     imageName,
		Cmd:       []string{"/bin/sh"},
		Tty:       true,
		OpenStdin: true,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}
	c.ContainerID = resp.ID
	// Start the container
	err = c.Start()
	if err != nil {
		return err
	}

	// update apk cache
	// log.Println("Updating apk cache")
	_, err = c.ExecCommand("apk update")
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) Start() error {
	err := c.DockerClient.ContainerStart(c.ctx, c.ContainerID, container.StartOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Stop stops (or destroys) the container
func (c *Container) Stop() error {
	err := c.DockerClient.ContainerRemove(c.ctx, c.ContainerID, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
		RemoveLinks:   false,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) Wait() error {
	statusCh, errCh := c.DockerClient.ContainerWait(c.ctx, c.ContainerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}
	return nil
}

func (c *Container) Logs() (string, error) {
	out, err := c.DockerClient.ContainerLogs(c.ctx, c.ContainerID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	output, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (c *Container) ExecCommand(command string) (string, error) {
	execConfig := types.ExecConfig{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := c.DockerClient.ContainerExecCreate(c.ctx, c.ContainerID, execConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.DockerClient.ContainerExecAttach(c.ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// pullImage pulls the specified Docker image
func pullImage(cli *client.Client, imageName string) error {
	ctx := context.Background()

	// dont pull the image if it already exists
	_, _, err := cli.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		// image exists, lets goo!
		return nil
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	return nil
}

func BuildImage(cli *client.Client, ctx context.Context, buildContext io.Reader, tags []string) error {
	opts := types.ImageBuildOptions{
		Tags: tags,
	}

	resp, err := cli.ImageBuild(ctx, buildContext, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
