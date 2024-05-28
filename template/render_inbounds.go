package template

import (
	"net/netip"

	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-dns"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json/badjson"
)

func (t *Template) renderInbounds(metadata M.Metadata, options *option.Options) error {
	options.Inbounds = t.Inbounds
	var needSniff bool
	if !t.DisableTrafficBypass {
		needSniff = true
	}
	var domainStrategy option.DomainStrategy
	if !t.RemoteResolve {
		if t.DomainStrategy != option.DomainStrategy(dns.DomainStrategyAsIS) {
			domainStrategy = t.DomainStrategy
		} else {
			domainStrategy = option.DomainStrategy(dns.DomainStrategyPreferIPv4)
		}
	}
	autoRedirect := t.AutoRedirect &&
		!metadata.Platform.IsApple() &&
		(metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.10.0-alpha.2")))
	disableTun := t.DisableTUN && !metadata.Platform.TunOnly()
	if !disableTun {
		options.Route.AutoDetectInterface = true

		var inet6Address []netip.Prefix
		if !t.DisableIPv6() {
			inet6Address = []netip.Prefix{netip.MustParsePrefix("fdfe:dcba:9876::1/126")}
		}
		tunInbound := option.Inbound{
			Type: C.TypeTun,
			TunOptions: option.TunInboundOptions{
				Inet4Address: []netip.Prefix{netip.MustParsePrefix("172.19.0.1/30")},
				Inet6Address: inet6Address,
				AutoRoute:    true,
				InboundOptions: option.InboundOptions{
					SniffEnabled: needSniff,
				},
			},
		}
		if autoRedirect {
			tunInbound.TunOptions.AutoRedirect = true
		}
		if t.EnableFakeIP {
			tunInbound.TunOptions.InboundOptions.DomainStrategy = domainStrategy
		}
		if metadata.Platform == M.PlatformUnknown {
			tunInbound.TunOptions.StrictRoute = true
		}
		if !t.DisableSystemProxy && metadata.Platform != M.PlatformUnknown {
			var httpPort uint16
			if t.CustomMixed != nil {
				httpPort = t.CustomMixed.Value.ListenPort
			}
			if httpPort == 0 {
				httpPort = DefaultMixedPort
			}
			tunInbound.TunOptions.Platform = &option.TunPlatformOptions{
				HTTPProxy: &option.HTTPProxyOptions{
					Enabled: true,
					ServerOptions: option.ServerOptions{
						Server:     "127.0.0.1",
						ServerPort: httpPort,
					},
				},
			}
		}
		if t.CustomTUN != nil {
			newTUNOptions, err := badjson.MergeFromDestination(tunInbound.TunOptions, t.CustomTUN.Message)
			if err != nil {
				return E.Cause(err, "merge custom tun options")
			}
			tunInbound.TunOptions = newTUNOptions
		}
		options.Inbounds = append(options.Inbounds, tunInbound)
	}
	if disableTun || !t.DisableSystemProxy {
		mixedInbound := option.Inbound{
			Type: C.TypeMixed,
			MixedOptions: option.HTTPMixedInboundOptions{
				ListenOptions: option.ListenOptions{
					Listen:     option.NewListenAddress(netip.AddrFrom4([4]byte{127, 0, 0, 1})),
					ListenPort: DefaultMixedPort,
					InboundOptions: option.InboundOptions{
						SniffEnabled:   needSniff,
						DomainStrategy: domainStrategy,
					},
				},
				SetSystemProxy: metadata.Platform == M.PlatformUnknown && disableTun && !t.DisableSystemProxy,
			},
		}
		if t.CustomMixed != nil {
			newMixedOptions, err := badjson.MergeFromDestination(mixedInbound.MixedOptions, t.CustomMixed.Message)
			if err != nil {
				return E.Cause(err, "merge custom mixed options")
			}
			mixedInbound.MixedOptions = newMixedOptions
		}
		options.Inbounds = append(options.Inbounds, mixedInbound)
	}
	return nil
}
