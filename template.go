package serenity

import (
	"net/netip"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-dns"
	N "github.com/sagernet/sing/common/network"
)

func DefaultTemplate(profileName string, platform string, version *Version, debug bool) *Profile {
	if version != nil && version.EqualAfter(ParseVersion("1.8.0-alpha.6")) {
		return defaultTemplate18a6(profileName, platform, version, debug)
	} else if version == nil || version.EqualAfter(ParseVersion("1.8.0-alpha.1")) {
		return defaultTemplate18(profileName, platform, version, debug)
	} else {
		return defaultTemplate17(profileName, platform, version, debug)
	}
}

func defaultTemplate18a6(profileName string, platform string, version *Version, debug bool) *Profile {
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
					Outbound: []string{"any"},
					Server:   "local",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					ClashMode: "Direct",
					Server:    "local",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					ClashMode: "Global",
					Server:    "google",
				},
			},
			{
				Type: C.RuleTypeLogical,
				LogicalOptions: option.LogicalDNSRule{
					Mode: C.LogicalTypeAnd,
					Rules: []option.DNSRule{
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultDNSRule{
								RuleSet: []string{"geosite-geolocation-!cn"},
								Invert:  true,
							},
						},
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultDNSRule{
								RuleSet: []string{
									"geosite-cn",
									"geosite-category-companies@cn",
								},
								DomainSuffix: []string{"download.jetbrains.com"},
							},
						},
					},
					Server: "local",
				},
			},
		},
	}
	options.Inbounds = []option.Inbound{
		{
			Type: C.TypeTun,
			TunOptions: option.TunInboundOptions{
				Inet4Address:           []netip.Prefix{netip.MustParsePrefix("172.19.0.1/30")},
				AutoRoute:              true,
				StrictRoute:            true,
				EndpointIndependentNat: true,
				UDPTimeout:             60,
				InboundOptions: option.InboundOptions{
					SniffEnabled: true,
				},
				Platform: &option.TunPlatformOptions{
					HTTPProxy: &option.HTTPProxyOptions{
						Enabled: true,
						ServerOptions: option.ServerOptions{
							Server:     "127.0.0.1",
							ServerPort: 8100,
						},
					},
				},
			},
		},
		{
			Type: C.TypeHTTP,
			HTTPOptions: option.HTTPMixedInboundOptions{
				ListenOptions: option.ListenOptions{
					Listen:     option.NewListenAddress(netip.AddrFrom4([4]byte{127, 0, 0, 1})),
					ListenPort: 8100,
					InboundOptions: option.InboundOptions{
						SniffEnabled:   true,
						DomainStrategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
					},
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
					IPIsPrivate: true,
					Outbound:    "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					ClashMode: "Direct",
					Outbound:  "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					ClashMode: "Global",
					Outbound:  "default",
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
				Type: C.RuleTypeLogical,
				LogicalOptions: option.LogicalRule{
					Mode: C.LogicalTypeAnd,
					Rules: []option.Rule{
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultRule{
								RuleSet: []string{"geosite-geolocation-!cn"},
								Invert:  true,
							},
						},
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultRule{
								RuleSet:      []string{"geoip-cn", "geosite-cn", "geosite-category-companies@cn"},
								DomainSuffix: []string{"download.jetbrains.com"},
							},
						},
					},
					Outbound: "direct",
				},
			},
		},
		RuleSet: []option.RuleSet{
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geoip-cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL: "https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs",
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL: "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs",
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-geolocation-!cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL: "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-geolocation-!cn.srs",
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-category-companies@cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL: "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-category-companies@cn.srs",
				},
			},
		},
		Final:               "default",
		AutoDetectInterface: true,
	}
	options.Experimental = &option.ExperimentalOptions{
		CacheFile: &option.CacheFileOptions{
			Enabled: true,
			CacheID: profileName,
		},
	}
	if platform == "" {
		options.Experimental.ClashAPI = &option.ClashAPIOptions{
			ExternalController: "127.0.0.1:9090",
			ExternalUI:         "clash-dashboard",
		}
	}
	return &Profile{
		options:  options,
		groupTag: []string{"default"},
	}
}

func defaultTemplate18(profileName string, platform string, version *Version, debug bool) *Profile {
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
					Outbound: []string{"any"},
					Server:   "local",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					ClashMode: "Direct",
					Server:    "local",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					ClashMode: "Global",
					Server:    "google",
				},
			},
			{
				Type: C.RuleTypeLogical,
				LogicalOptions: option.LogicalDNSRule{
					Mode: C.LogicalTypeAnd,
					Rules: []option.DNSRule{
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultDNSRule{
								RuleSet: []string{"geosite-geolocation-!cn"},
								Invert:  true,
							},
						},
						{
							Type: C.RuleTypeLogical,
							LogicalOptions: option.LogicalDNSRule{
								Mode: C.LogicalTypeOr,
								Rules: []option.DNSRule{
									{
										Type: C.RuleTypeDefault,
										DefaultOptions: option.DefaultDNSRule{
											RuleSet: []string{
												"geosite-cn",
												"geosite-category-companies@cn",
											},
										},
									},
									{
										Type: C.RuleTypeDefault,
										DefaultOptions: option.DefaultDNSRule{
											DomainSuffix: []string{"download.jetbrains.com"},
										},
									},
								},
							},
						},
					},
					Server: "local",
				},
			},
		},
	}
	options.Inbounds = []option.Inbound{
		{
			Type: C.TypeTun,
			TunOptions: option.TunInboundOptions{
				Inet4Address:           []netip.Prefix{netip.MustParsePrefix("172.19.0.1/30")},
				AutoRoute:              true,
				StrictRoute:            true,
				EndpointIndependentNat: true,
				UDPTimeout:             60,
				InboundOptions: option.InboundOptions{
					SniffEnabled: true,
				},
				Platform: &option.TunPlatformOptions{
					HTTPProxy: &option.HTTPProxyOptions{
						Enabled: true,
						ServerOptions: option.ServerOptions{
							Server:     "127.0.0.1",
							ServerPort: 8100,
						},
					},
				},
			},
		},
		{
			Type: C.TypeHTTP,
			HTTPOptions: option.HTTPMixedInboundOptions{
				ListenOptions: option.ListenOptions{
					Listen:     option.NewListenAddress(netip.AddrFrom4([4]byte{127, 0, 0, 1})),
					ListenPort: 8100,
					InboundOptions: option.InboundOptions{
						SniffEnabled:   true,
						DomainStrategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
					},
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
					GeoIP:    []string{"private"},
					Outbound: "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					ClashMode: "Direct",
					Outbound:  "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					ClashMode: "Global",
					Outbound:  "default",
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
				Type: C.RuleTypeLogical,
				LogicalOptions: option.LogicalRule{
					Mode: C.LogicalTypeAnd,
					Rules: []option.Rule{
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultRule{
								RuleSet: []string{"geosite-geolocation-!cn"},
								Invert:  true,
							},
						},
						{
							Type: C.RuleTypeLogical,
							LogicalOptions: option.LogicalRule{
								Mode: C.LogicalTypeOr,
								Rules: []option.Rule{
									{
										Type: C.RuleTypeDefault,
										DefaultOptions: option.DefaultRule{
											RuleSet: []string{"geoip-cn", "geosite-cn", "geosite-category-companies@cn"},
										},
									},
									{
										Type: C.RuleTypeDefault,
										DefaultOptions: option.DefaultRule{
											DomainSuffix: []string{"download.jetbrains.com"},
										},
									},
								},
							},
						},
					},
					Outbound: "direct",
				},
			},
		},
		RuleSet: []option.RuleSet{
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geoip-cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL: "https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs",
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL: "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs",
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-geolocation-!cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL: "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-geolocation-!cn.srs",
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-category-companies@cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL: "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-category-companies@cn.srs",
				},
			},
		},
		Final:               "default",
		AutoDetectInterface: true,
	}
	options.Experimental = &option.ExperimentalOptions{
		CacheFile: &option.CacheFileOptions{
			Enabled: true,
			CacheID: profileName,
		},
	}
	if platform == "" {
		options.Experimental.ClashAPI = &option.ClashAPIOptions{
			ExternalController: "127.0.0.1:9090",
			ExternalUI:         "clash-dashboard",
		}
	}
	return &Profile{
		options:  options,
		groupTag: []string{"default"},
	}
}

//goland:noinspection GoDeprecation
//nolint:staticcheck
func defaultTemplate17(profileName string, platform string, version *Version, debug bool) *Profile {
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
					Outbound: []string{"any"},
					Server:   "local",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					ClashMode: "Direct",
					Server:    "local",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultDNSRule{
					ClashMode: "Global",
					Server:    "google",
				},
			},
			{
				Type: C.RuleTypeLogical,
				LogicalOptions: option.LogicalDNSRule{
					Mode: C.LogicalTypeAnd,
					Rules: []option.DNSRule{
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultDNSRule{
								Geosite: []string{"geolocation-!cn"},
								Invert:  true,
							},
						},
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultDNSRule{
								Geosite: []string{
									"cn",
									"category-companies@cn",
								},
							},
						},
					},
					Server: "local",
				},
			},
		},
	}
	options.Inbounds = []option.Inbound{
		{
			Type: C.TypeTun,
			TunOptions: option.TunInboundOptions{
				Inet4Address:           []netip.Prefix{netip.MustParsePrefix("172.19.0.1/30")},
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
					GeoIP:    []string{"private"},
					Outbound: "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					ClashMode: "Direct",
					Outbound:  "direct",
				},
			},
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					ClashMode: "Global",
					Outbound:  "default",
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
				Type: C.RuleTypeLogical,
				LogicalOptions: option.LogicalRule{
					Mode: C.LogicalTypeAnd,
					Rules: []option.Rule{
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultRule{
								Geosite: []string{"geolocation-!cn"},
								Invert:  true,
							},
						},
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultRule{
								GeoIP:        []string{"cn"},
								Geosite:      []string{"cn", "category-companies@cn"},
								DomainSuffix: []string{"download.jetbrains.com"},
							},
						},
					},
					Outbound: "direct",
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
	if debug && (version == nil || version.After(ParseVersion("1.3-beta8"))) {
		options.Experimental.Debug = &option.DebugOptions{
			Listen: "0.0.0.0:8965",
		}
	}
	if version == nil || version.After(ParseVersion("1.3-beta1")) {
		options.Experimental.ClashAPI.ExternalUI = "clash-dashboard"
	}
	if version == nil || version.After(ParseVersion("1.3-beta11")) {
		options.Experimental.ClashAPI.CacheID = profileName
	}
	if version == nil || version.After(ParseVersion("1.4.0-beta.3")) {
		options.Inbounds[0].TunOptions.Platform = &option.TunPlatformOptions{
			HTTPProxy: &option.HTTPProxyOptions{
				Enabled: true,
				ServerOptions: option.ServerOptions{
					Server:     "127.0.0.1",
					ServerPort: 8100,
				},
			},
		}
		options.Inbounds = append(options.Inbounds, option.Inbound{
			Type: C.TypeHTTP,
			HTTPOptions: option.HTTPMixedInboundOptions{
				ListenOptions: option.ListenOptions{
					Listen:     option.NewListenAddress(netip.AddrFrom4([4]byte{127, 0, 0, 1})),
					ListenPort: 8100,
					InboundOptions: option.InboundOptions{
						SniffEnabled:   true,
						DomainStrategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
					},
				},
			},
		})
	}
	if version == nil || version.After(ParseVersion("1.4.0-rc.2")) {
		options.Experimental.ClashAPI.StoreMode = true
	}
	return &Profile{
		options:  options,
		groupTag: []string{"default"},
	}
}
