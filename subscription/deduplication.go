package subscription

import (
	"context"
	"net/netip"
	"sync"

	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-dns"
	"github.com/sagernet/sing/common"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/common/task"
)

func Deduplication(ctx context.Context, servers []option.Outbound) []option.Outbound {
	resolveCtx := &resolveContext{
		ctx: ctx,
		dnsClient: dns.NewClient(dns.ClientOptions{
			DisableExpire: true,
			Logger:        log.NewNOPFactory().Logger(),
		}),
		dnsTransport: common.Must1(dns.NewTLSTransport(dns.TransportOptions{
			Context:      ctx,
			Dialer:       N.SystemDialer,
			Address:      "tls://1.1.1.1",
			ClientSubnet: netip.MustParseAddr("114.114.114.114"),
		})),
	}

	uniqueServers := make([]netip.AddrPort, len(servers))
	var (
		resolveGroup task.Group
		resultAccess sync.Mutex
	)
	for index, server := range servers {
		currentIndex := index
		currentServer := server
		resolveGroup.Append0(func(ctx context.Context) error {
			destination := resolveDestination(resolveCtx, currentServer)
			if destination.IsValid() {
				resultAccess.Lock()
				uniqueServers[currentIndex] = destination
				resultAccess.Unlock()
			}
			return nil
		})
		resolveGroup.Concurrency(5)
		_ = resolveGroup.Run(ctx)
	}
	uniqueServerMap := make(map[netip.AddrPort]bool)
	var newServers []option.Outbound
	for index, server := range servers {
		destination := uniqueServers[index]
		if destination.IsValid() {
			if uniqueServerMap[destination] {
				continue
			}
			uniqueServerMap[destination] = true
		}
		newServers = append(newServers, server)
	}
	return newServers
}

type resolveContext struct {
	ctx          context.Context
	dnsClient    *dns.Client
	dnsTransport dns.Transport
}

func resolveDestination(ctx *resolveContext, server option.Outbound) netip.AddrPort {
	rawOptions, err := server.RawOptions()
	if err != nil {
		return netip.AddrPort{}
	}
	serverOptionsWrapper, loaded := rawOptions.(option.ServerOptionsWrapper)
	if !loaded {
		return netip.AddrPort{}
	}
	serverOptions := serverOptionsWrapper.TakeServerOptions().Build()
	if serverOptions.IsIP() {
		return serverOptions.AddrPort()
	}
	if serverOptions.IsFqdn() {
		addresses, lookupErr := ctx.dnsClient.Lookup(ctx.ctx, ctx.dnsTransport, serverOptions.Fqdn, dns.DomainStrategyPreferIPv4)
		if lookupErr == nil && len(addresses) > 0 {
			return netip.AddrPortFrom(addresses[0], serverOptions.Port)
		}
	}
	return netip.AddrPort{}
}
