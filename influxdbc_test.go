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
