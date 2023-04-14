package serenity

import (
	"net/netip"

	"github.com/sagernet/sing-box/common/badversion"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-dns"
	N "github.com/sagernet/sing/common/network"
)

func DefaultTemplate(platform string, version *badversion.Version) *Profile {
	var options option.Options
	options.Log = &option.LogOptions{
		Level: "info",
	}
	options.DNS = &option.DNSOptions{
		DNSClientOptions: option.DNSClientOptions{
			Strategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
		},
		Servers: []option.DNSServerOptions{
			{
				Tag:     "google",
				Address: "tls://8.8.8.8",
			},
			{
				Tag:     "local",
				Address: "114.114.114.114",
				Detour:  "direct",
			},
		},
		Rules: []option.DNSRule{
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					ClashMode: "direct",
					Server:    "local",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					Geosite: []string{"cn"},
					Server:  "local",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					Outbound: []string{"any"},
					Server:   "local",
				},
			},
		},
	}
	options.Inbounds = []option.Inbound{
		{
			Type: C.TypeTun,
			TunOptions: option.TunInboundOptions{
				Inet4Address:           []option.ListenPrefix{option.ListenPrefix(netip.MustParsePrefix("172.19.0.1/30"))},
				AutoRoute:              true,
				StrictRoute:            true,
				EndpointIndependentNat: true,
				UDPTimeout:             60,
				InboundOptions: option.InboundOptions{
					SniffEnabled: true,
				},
			},
		},
	}
	options.Outbounds = []option.Outbound{
		{
			Tag:  "default",
			Type: C.TypeSelector,
		},
		{
			Tag:  "direct",
			Type: C.TypeDirect,
		},
		{
			Tag:  "block",
			Type: C.TypeBlock,
		},
		{
			Tag:  "dns",
			Type: C.TypeDNS,
		},
	}
	options.Route = &option.RouteOptions{
		Rules: []option.Rule{
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					Protocol: []string{"dns"},
					Outbound: "dns",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					Network:  []string{N.NetworkUDP},
					Port:     []uint16{53},
					Outbound: "dns",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					Port:     []uint16{853},
					Outbound: "block",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					Network:  []string{N.NetworkUDP},
					Port:     []uint16{443},
					Outbound: "block",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					ClashMode: "direct",
					Outbound:  "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					Protocol: []string{"stun"},
					Outbound: "block",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					GeoIP:        []string{"cn", "private"},
					Geosite:      []string{"cn"},
					DomainSuffix: []string{"download.jetbrains.com"},
					Outbound:     "direct",
				},
			},
		},
		Final:               "default",
		AutoDetectInterface: true,
	}
	options.Experimental = &option.ExperimentalOptions{
		ClashAPI: &option.ClashAPIOptions{
			ExternalController: "127.0.0.1:9090",
			StoreSelected:      true,
		},
	}
	if version == nil || version.After(badversion.Parse("1.3-beta1")) {
		options.Experimental.ClashAPI.ExternalUI = "clash-dashboard"
	}
	return &Profile{
		options:  options,
		groupTag: []string{"default"},
	}
}
