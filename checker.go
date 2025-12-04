package main

import "context"

type Checker interface {
	Check(ctx context.Context, ep Endpoint) (Result, error)
}

type TCPChecker struct{}

func (c *TCPChecker) Check(ctx context.Context, ep Endpoint) (Result, error) {
	panic("TODO")
}

type DNSChecker struct{}

// etc.
