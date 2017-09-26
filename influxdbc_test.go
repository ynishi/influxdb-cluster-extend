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
	ctx context.Context
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
	node, _ = NewNode(ctx)
	masterNode, _ = NewMasterNode(ctx)
	cluster, _ = NewCluster(ctx)
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
		host,
		port,
		ctx,
	}

	if !reflect.DeepEqual(node, testNode) {
		t.Fatalf("Not match node:\n want: %q,\n have: %q\n",node, testNode)
	}
}

func TestNewNode(t *testing.T) {

	testNode := NewNode(ctx)
	if !reflect.DeepEqual(node, testNode) {
		t.Fatalf("Not match node by NewNode:\n want: %q,\n have: %q\n",node, testNode)
	}
}

func TestCheckNode(t *testing.T) {

	res, err := node.Check()
	if err != nil {
		t.Fatal(err)
	}
	if res != OK {
		t.Fatalf("Check failed:\n want: %q,\n have: %q\n", OK, res)
	}
}

func TestAddNode(t *testing.T) {

	testCluster := NewCluster(ctx)
	nodeId, err := testCluster.Add(node)
	if err != nil {
		t.Fatal(err)
	}
	if nodeId != id {
		t.Fatal("Not match node id:\n want:%q,\n have: %q\n",id,nodeId)
	}
	if !reflect.DeepEqual(cluster, testCluster) {
		t.Fatalf("Not match nodes:\n want:%q,\n have: %q\n",nodes,testNodes)
	}
}

func TestNodeIds(t *testing.T) {

	testIds, err := cluster.Ids()
	if err != nil {
		t.Fatal(err)
	}
	if testIds != ids {
		t.Fatal("Not match node ids:\n want:%q,\n have: %q\n",ids,testIds)
	}
}

func TestRemoveNode(t *testing.T) {

	testCluster := NewCluster(ctx)
	nodeId, err := testCluster.Add(node)
	err := cluster.Remove(nodeId)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range cluster.Ids() {
		if id == nodeId {
			t.Fatalf("Failed remove node: %q remains\n", nodeId)
		}
	}

}

func TestNewMasterNode(t *testing.T){

	testMasterNode, err := NewMasterNode(ctx)
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
	if res != OK {
		t.Fatalf("Master Node check failed:\n want: %q,\n have: %q\n", OK, res)
	}
}

func TestCountNodeCluster(t *testing.T) {

	if cluster.Count() != 1 {
		t.Fatalf("Not match testNodes count:\n want: %q,\n have:%q\n", 1, cluster.Count())
	}

}

func TestStartCluster(t *testing.T) {

	testMasterNode := NewMasterNode(ctx, &MasterNodeConfig{})
	testCluster := NewCluster(ctx, testMaterNode)
	err := testCluster.Start()
	if err != nil {
		t.Fatal(err)
	}
}


func TestFailOver(t *testing.T) {

	testMasterNode := cluster.GetMasterNode()

	testNewMasterId, err := cluster.FailOver()
	if err != nil {
		t.Fatal(err)
	}
	if testNewMasterId == testMasterNode.Id {
		t.Fatal("failed failover, master node not changed.")
	}
}

func TestCheckCluster