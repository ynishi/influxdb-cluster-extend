package influxdbc

import (
	"context"
	"reflect"
	"testing"
	"golang.org/x/exp/shiny/widget/node"
	"github.com/docker/docker/pkg/discovery/nodes"
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

func TestCheckNode(t *testing.T) {

	testNode := NewNode(ctx)
	res, err := testNode.Check()
	if err != nil {
		t.Fatal(err)
	}
	if res != Ok {
		t.Fatalf("Check failed:\n want: %q,\n have: %q\n",Ok, res)
	}
}

func TestCountNode(t *testing.T) {

	if testNodes.Count() != 1 {
		t.Fatalf("Not match testNodes count:\n want: %q,\n have:%q\n", 1, testNodes.Count())
	}

}

func TestAddNode(t *testing.T) {

	testNode := NewNode(ctx)
	_, err := testNodes.add(testNode)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(nodes, testNodes) {
		t.Fatalf("Not match nodes:\n want:%q,\n have: %q\n",nodes,testNodes)
	}
}

func TestNodeId(t *testing.T) {

	testNode := NewNode(ctx)
	testId, err := testNodes.add(testNode)
	if err != nil {
		t.Fatal(err)
	}
	if testId != id {
		t.Fatal("Not match node id:\n want:%q,\n have: %q\n",id,testId)
	}
}

func TestRemoveNode(t *testing.T) {

	testNode := NewNode(ctx)
	testId, err := testNodes.add(testNode)
	if err != nil {
		t.Fatal(err)
	}
	err := testNodes.remove(testId)
	if err != nil {
		t.Fatal(err)
	}
	if testNodes.count() != 0 {
		t.Fatalf("Node count > 0:\n want: 0,\n have: %q\n", testNodes.count())
	}
}