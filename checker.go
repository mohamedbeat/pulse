package main

import "context"

type Checker interface {
	Check(ctx context.Context, ep Endpoint) Result
}

type TCPChecker struct{}

func (c *TCPChecker) Check(ctx context.Context, ep Endpoint) Result {
	panic("TODO")
}

type DNSChecker struct{}

// etc.
