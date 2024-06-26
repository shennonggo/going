package direct

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

type directBuilder struct{}

func init() {
	resolver.Register(NewBuilder())
}

/*
NewBuilder creates a directBuilder which is used to factory direct resolvers.
example:

	direct://<authority>/127.0.0.1:9000

到达这样一个目的
*/
func NewBuilder() *directBuilder {
	return &directBuilder{}
}
func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, 0)
	for _, addr := range strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}
	//grpc建立连接的逻辑都在这里UpdateState
	err := cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		return nil, err
	}
	return newDirectResolver(), nil
}

func (d *directBuilder) Scheme() string {
	return "direct"
}

var _ resolver.Builder = &directBuilder{}
