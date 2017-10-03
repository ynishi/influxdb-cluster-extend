package influxdbc

import (
	"context"
	"errors"
	"log"
	"time"

	"fmt"

	eclient "github.com/coreos/etcd/client"
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

type Node struct {
	Checker
	StartStopper
	Host string
	Port string
	Etcd eclient.Client
	Apic eclient.KeysAPI
	Ctx  context.Context
}

type MasterNode struct {
	Checker
	Host string
	Port string
	Etcd eclient.Client
	Apic eclient.KeysAPI
	Ctx  context.Context
}

type Cluster struct {
	Checker
	StartStopper
	Endpoint     string
	Port         string
	Nodes        map[string]*Node
	MasterNode   *MasterNode
	MasterNodeId string
	Ctx          context.Context
}

type NodeConfig struct {
	Host string
	Port string
}

type ClusterConfig struct {
	Endpoint string
	Port     string
}

type Checker interface {
	Check() *CheckResult
}

type Starter interface {
	Start() error
}

type Stopper interface {
	Stop() error
}

type StartStopper interface {
	Starter
	Stopper
}

func NewNode(ctx context.Context, nodeConfig *NodeConfig) (node *Node, err error) {

	if nodeConfig == nil {
		return nil, errors.New("node config is nil.")
	}
	cfg := eclient.Config{
		Endpoints:               []string{fmt.Sprintf("%s:2379", nodeConfig.Host)},
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
		node = &Node{Ctx: ctx}
	}
	return node, nil
}

func (n *Node) Start() (err error) {
	return nil
}

func (n *Node) stop() (err error) {
	return nil
}

func (n *Node) Check() (res *CheckResult, err error) {
	resp, err := n.Apic.Get(n.Ctx, "/status", nil)
	if err != nil {
		return nil, err
	}
	if resp.Node.Value == "CRITICAL" {
		return &CheckResult{
			CRITICAL,
			"CRITICAL",
		}, nil
	} else if resp.Node.Value == "WARN" {
		return &CheckResult{
			WARN,
			"WARN",
		}, nil
	} else if resp.Node.Value == "OK" {
		return &CheckResult{
			OK,
			"OK",
		}, nil
	} else {
		return nil, nil
	}
}

func NewMasterNode(ctx context.Context, nodeConfig *NodeConfig) (node *MasterNode, err error) {

	if nodeConfig == nil {
		return nil, errors.New("node config is nil.")
	}
	cfg := eclient.Config{
		Endpoints:               []string{fmt.Sprintf("%s:2379", nodeConfig.Host)},
		Transport:               eclient.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := eclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := eclient.NewKeysAPI(c)
	if nodeConfig != nil {
		node = &MasterNode{
			Host: nodeConfig.Host,
			Port: nodeConfig.Port,
			Ctx:  ctx,
			Etcd: c,
			Apic: kapi,
		}
	} else {
		node = &MasterNode{Ctx: ctx}
	}
	return node, nil
}

func (mn *MasterNode) Check() (res *CheckResult, err error) {
	resp, err := mn.Apic.Get(mn.Ctx, "/status", nil)
	if err != nil {
		return nil, err
	}
	if resp.Node.Value == "CRITICAL" {
		return &CheckResult{
			CRITICAL,
			"CRITICAL",
		}, nil
	} else if resp.Node.Value == "OK" {
		return &CheckResult{
			OK,
			"OK",
		}, nil
	} else {
		return nil, nil
	}
}

func NewCluster(ctx context.Context, clusterConfig *ClusterConfig) (cluster *Cluster, err error) {

	if clusterConfig != nil {
		cluster = &Cluster{
			Endpoint: clusterConfig.Endpoint,
			Port:     clusterConfig.Port,
			Ctx:      ctx,
		}
	} else {
		cluster = &Cluster{Ctx: ctx}
	}
	return cluster, nil
}

func (c *Cluster) Start() (err error) {

	for _, node := range c.Nodes {
		err := node.Start()
		if err != nil {
			return err
		}
	}
	_, err = c.Check()
	if err != nil {
		c.Stop()
		return errors.New("failed start cluster: "+err.Error())
	}
	return nil
}

func (c *Cluster) Stop() (err error) {

	isError := false
	for id, node := range c.Nodes {
		err := node.Stop()
		if err != nil {
			isError = true
			log.Println("failed stop node: %s, %v",id,err)
		}
	}
	if isError {
		return errors.New("failed stop cluster: "+err.Error())
	} else {
		return nil
	}
}

func (c *Cluster) Add(node *Node) (id string, err error) {
	if node == nil {
		return "", errors.New("cannot add node, node is nil")
	}
	id = mkHashId(node)
	c.Nodes[id] = node
	return id, nil
}

func (c *Cluster) Remove(id string) (err error) {
	return nil
}

func mkHashId(node *Node) (id string) {
	return ""
}

func (c *Cluster) FailOver() (err error) {
	return nil
}

func (c *Cluster) Check() (res *CheckResult, err error) {
	var checkOk, checkNG []CheckResult
	for _, n := range c.Nodes {
		res, err := n.Check()
		if err != nil {
			return nil, err
		}
		if res.Code == OK {
			checkOk = append(checkOk, *res)
		} else {
			checkNG = append(checkNG, *res)
		}
	}

	if len(checkNG) > len(checkOk) {
		return &CheckResult{
			CRITICAL,
			"CRITICAL",
		}, nil
	} else {
		return &CheckResult{
			OK,
			"OK",
		}, nil
	}
}
