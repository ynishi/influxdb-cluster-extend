package influxdbc

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dclient "github.com/docker/docker/client"
	"github.com/influxdata/influxdb"
	eclient "github.com/coreos/etcd/client"
	"fmt"
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

type Node struct {
	Checker
	Host string
	Port string
	Etcd eclient.Node
	Apic eclient.Kapi
	Ctx context.Context
}

type NodeConfig struct {
	Host string
	Port string
}

type Checker interface {
	Check() *CheckResult
}

func NewNode(ctx context.Context, nodeConfig *NodeConfig) (node *Node, err error) {

	cfg := eclient.Config{
		Endpoints:               []string{fmt.Sprintf("%s:2379",nodeConfig.Host)},
		Transport:               eclient.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := eclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := eclient.NewKeysAPI(c)
	if nodeConfig != nil {
		node = &Node{
			Host: nodeConfig.Host,
			Port: nodeConfig.Port,
			Ctx:  ctx,
			Etcd: c,
			Apic: kapi,
		}
	} else {
		node = &Node{Ctx:ctx}
	}
	return node, nil
}

func (n *Node) Check() (res *CheckResult,err error) {
	resp, err := n.Apic.Get(n.Ctx, "/status", nil)
	if err != nil {
		return nil, err
	}
	if resp == "CRITICAL" {
		return &CheckResult{
			CRITICAL,
			"CRITICAL",
		}, nil
	} else if resp == "WARN"{
		return &CheckResult{
			WARN,
			"WARN",
		},nil
	} else if resp == "OK" {
		return &CheckResult{
			OK,
			"OK",
		}, nil
	} else {
		return nil, nil
	}
}