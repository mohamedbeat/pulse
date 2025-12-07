package main

import (
	"context"

	"github.com/mohamedbeat/pulse/common"
)

type Checker interface {
	Check(ctx context.Context, ep common.Endpoint) common.Result
}

type TCPChecker struct{}

func (c *TCPChecker) Check(ctx context.Context, ep common.Endpoint) common.Result {
	panic("TODO")
}

type DNSChecker struct{}

// etc.
