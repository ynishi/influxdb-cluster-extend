package influxdbc

import (
	"context"
	"reflect"
	"testing"
)

var (
	client      *Client
	checkResult *CheckResult
	idResult    string
	ctx         context.Context
	node        *Node
	masterNode  *MasterNode
	cluster     *Cluster
	id          string
)

func init() {

	checkResult = &CheckResult{
		Code:    OK,
		Message: "OK",
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, _ = NewClient(ctx)
	node, _ = NewNode(ctx, nil)
	masterNode, _ = NewMasterNode(ctx, nil)
	cluster, _ = NewCluster(ctx, nil)
	id = ""
}

func TestRun(t *testing.T) {

	errChan := make(chan error, 1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		errChan <- Run(ctx)
	}()
	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	case err := <-errChan:
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCheck(t *testing.T) {

	resultChan := make(chan *CheckResult, 1)
	errChan := make(chan error, 1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		result, err := Check(ctx)
		if err != nil {
			errChan <- err
		}
		resultChan <- result
	}()
	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	case err := <-errChan:
		if err != nil {
			t.Fatal(err)
		}
	case testResult := <-resultChan:
		if !reflect.DeepEqual(testResult, checkResult) {
			t.Fatal("failed match check.\n have: %q,\n want: %q\n", testResult, checkResult)
			t.Fatal("failed match check.\n have: %q,\n want: %q\n", testResult, checkResult)
		}
	}
}

func TestGetId(t *testing.T) {

	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		result, err := GetId(ctx)
		if err != nil {
			errChan <- err
		}
		resultChan <- result
	}()
	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	case err := <-errChan:
		if err != nil {
			t.Fatal(err)
		}
	case testResult := <-resultChan:
		if testResult != idResult {
			t.Fatal("failed match id.\n have: %v,\n want: %v\n", testResult, idResult)
		}
	}
}

func TestNewClient(t *testing.T) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	testClient, err := NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if testClient.ctx != ctx {
		t.Fatalf("Failed new client:\n want: %q,\n have: %q\n")
	}

}

func TestClient(t *testing.T) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	testClient := &Client{
		client.container,
		client.ctx,
	}

	if !reflect.DeepEqual(client, testClient) {
		t.Fatalf("Not match client:\n want: %q,\n have: %q\n", client, testClient)
	}
}

func TestNode(t *testing.T) {

	host := "127.0.0.1"
	port := "3999"
	testNode := &Node{
		Host: host,
		Port: port,
		Ctx:  ctx,
		Etcd: nil,
		Apic: nil,
	}

	if !reflect.DeepEqual(node, testNode) {
		t.Fatalf("Not match node:\n want: %q,\n have: %q\n", node, testNode)
	}
}

func TestNewNode(t *testing.T) {

	testNode, err := NewNode(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(node, testNode) {
		t.Fatalf("Not match node by NewNode:\n want: %q,\n have: %q\n", node, testNode)
	}
}

func TestCheckNode(t *testing.T) {

	res, err := node.Check()
	if err != nil {
		t.Fatal(err)
	}
	if res.Code != OK {
		t.Fatalf("Check failed:\n want: %q,\n have: %q\n", OK, res)
	}
}

func TestAddNode(t *testing.T) {

	testCluster, err := NewCluster(ctx, nil)
	nodeId, err := testCluster.Add(node)
	if err != nil {
		t.Fatal(err)
	}
	if nodeId != id {
		t.Fatal("Not match node id:\n want:%q,\n have: %q\n", id, nodeId)
	}
	if !reflect.DeepEqual(cluster, testCluster) {
		t.Fatalf("Not match nodes:\n want:%q,\n have: %q\n", cluster, testCluster)
	}
}

func TestRemoveNode(t *testing.T) {

	testCluster, err := NewCluster(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	nodeId, err := testCluster.Add(node)
	err = cluster.Remove(nodeId)
	if err != nil {
		t.Fatal(err)
	}
	for id, _ := range cluster.Nodes {
		if id == nodeId {
			t.Fatalf("Failed remove node: %q remains\n", nodeId)
		}
	}

}

func TestStopCluster(t *testing.T) {

	testCluster, err := NewCluster(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = testCluster.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewMasterNode(t *testing.T) {

	testMasterNode, err := NewMasterNode(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(testMasterNode, masterNode) {
		t.Fatalf("Failed new node:\n want: %q,\n have: %q\n", masterNode, testMasterNode)
	}
}

func TestCheckMasterNode(t *testing.T) {

	res, err := masterNode.Check()
	if err != nil {
		t.Fatal(err)
	}
	if res.Code != OK {
		t.Fatalf("Master Node check failed:\n want: %q,\n have: %q\n", OK, res)
	}
}

func TestStartCluster(t *testing.T) {

	testCluster, err := NewCluster(ctx, nil)
	err = testCluster.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFailOver(t *testing.T) {

	testMasterNodeId := cluster.MasterNodeId

	err := cluster.FailOver()
	if err != nil {
		t.Fatal(err)
	}
	if cluster.MasterNodeId == testMasterNodeId {
		t.Fatal("failed failover, master node not changed.")
	}
}
