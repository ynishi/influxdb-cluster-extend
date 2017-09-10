package influxdbc

import "context"

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

func NewClient(ctx context.Context) (client Client, err error) {
	return client, err
}
