package influxdbc

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dclient "github.com/docker/docker/client"
)

const (
	OK = iota
	WARN
	ERROR
	CRITICAL
)

type CheckResult struct {
	Code    int
	Message string
}

type Client struct {
	container container.ContainerCreateCreatedBody
	ctx       context.Context
}

var cli *dclient.Client

func init() {
	var err error
	cli, err = dclient.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}
}

func Run(ctx context.Context) (err error) {
	return nil
}

func Check(ctx context.Context) (check *CheckResult, err error) {
	check = &CheckResult{
		OK,
		"OK",
	}
	return check, err
}

func GetId(ctx context.Context) (id string, err error) {
	return id, err
}

func NewClient(ctx context.Context) (client *Client, err error) {

	containerConfig := &container.Config{
		Image: "influxdb",
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, "")
	if err != nil {
		return nil, err
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	client = &Client{
		container: resp,
		ctx:       ctx,
	}
	return client, err
}
