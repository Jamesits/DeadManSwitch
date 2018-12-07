package main

import (
	"context"
	"net"
)

// force using a specific server for DNS requests
func dialContextFactory(server string) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{DualStack: true}).DialContext(ctx, network, server)
	}
}

func getResolver() *net.Resolver {
	resolver := &net.Resolver{
		PreferGo: true,
	}

	return resolver
}

func getResolverWithServer(server string) *net.Resolver {
	resolver := &net.Resolver{
		PreferGo: true,
	}

	resolver.Dial = dialContextFactory(server)
	return resolver
}
