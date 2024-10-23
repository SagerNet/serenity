package template

import (
	"net/netip"
	"net/url"

	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-dns"
	"github.com/sagernet/sing/common"
	BM "github.com/sagernet/sing/common/metadata"

	mDNS "github.com/miekg/dns"
)

func (t *Template) renderDNS(metadata M.Metadata, options *option.Options) error {
	var (
		domainStrategy      option.DomainStrategy
		domainStrategyLocal option.DomainStrategy
	)
	if t.DomainStrategy != option.DomainStrategy(dns.DomainStrategyAsIS) {
		domainStrategy = t.DomainStrategy
	} else if t.EnableFakeIP {
		domainStrategy = option.DomainStrategy(dns.DomainStrategyPreferIPv4)
	} else {
		domainStrategy = option.DomainStrategy(dns.DomainStrategyUseIPv4)
	}
	if t.DomainStrategyLocal != option.DomainStrategy(dns.DomainStrategyAsIS) {
		domainStrategyLocal = t.DomainStrategyLocal
	} else {
		domainStrategyLocal = option.DomainStrategy(dns.DomainStrategyPreferIPv4)
	}
	if domainStrategyLocal == domainStrategy {
		domainStrategyLocal = 0
	}
	options.DNS = &option.DNSOptions{
		Servers:        t.Servers,
		ReverseMapping: !t.DisableTrafficBypass && metadata.Platform != M.PlatformUnknown && !metadata.Platform.IsApple(),
		DNSClientOptions: option.DNSClientOptions{
			Strategy:         domainStrategy,
			IndependentCache: t.EnableFakeIP,
		},
	}
	dnsDefault := t.DNS
	if dnsDefault == "" {
		dnsDefault = DefaultDNS
	}
	dnsLocal := t.DNSLocal
	if dnsLocal == "" {
		dnsLocal = DefaultDNSLocal
	}
	directTag := t.DirectTag
	if directTag == "" {
		directTag = DefaultDirectTag
	}
	defaultDNSOptions := option.DNSServerOptions{
		Tag:     DNSDefaultTag,
		Address: dnsDefault,
	}
	if dnsDefaultUrl, err := url.Parse(dnsDefault); err == nil && BM.IsDomainName(dnsDefaultUrl.Hostname()) {
		defaultDNSOptions.AddressResolver = DNSLocalTag
	}
	options.DNS.Servers = append(options.DNS.Servers, defaultDNSOptions)
	var (
		localDNSOptions  option.DNSServerOptions
		localDNSIsDomain bool
	)
	if t.DisableTrafficBypass {
		localDNSOptions = option.DNSServerOptions{
			Tag:      DNSLocalTag,
			Address:  "local",
			Strategy: domainStrategyLocal,
		}
	} else {
		localDNSOptions = option.DNSServerOptions{
			Tag:      DNSLocalTag,
			Address:  dnsLocal,
			Detour:   directTag,
			Strategy: domainStrategyLocal,
		}
		if dnsLocalUrl, err := url.Parse(dnsLocal); err == nil && BM.IsDomainName(dnsLocalUrl.Hostname()) {
			localDNSOptions.AddressResolver = DNSLocalSetupTag
			localDNSIsDomain = true
		}
	}
	options.DNS.Servers = append(options.DNS.Servers, localDNSOptions)
	if localDNSIsDomain {
		options.DNS.Servers = append(options.DNS.Servers, option.DNSServerOptions{
			Tag:      DNSLocalSetupTag,
			Address:  "local",
			Strategy: domainStrategyLocal,
		})
	}
	if t.EnableFakeIP {
		options.DNS.FakeIP = &option.DNSFakeIPOptions{
			Enabled:    true,
			Inet4Range: common.Ptr(netip.MustParsePrefix("198.18.0.0/15")),
		}
		if !t.DisableIPv6() {
			options.DNS.FakeIP.Inet6Range = common.Ptr(netip.MustParsePrefix("fc00::/18"))
		}
		options.DNS.Servers = append(options.DNS.Servers, option.DNSServerOptions{
			Tag:     DNSFakeIPTag,
			Address: "fakeip",
		})
	}
	options.DNS.Rules = []option.DNSRule{
		{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultDNSRule{
				Outbound: []string{"any"},
				Server:   DNSLocalTag,
			},
		},
	}
	clashModeRule := t.ClashModeRule
	if clashModeRule == "" {
		clashModeRule = "Rule"
	}
	clashModeGlobal := t.ClashModeGlobal
	if clashModeGlobal == "" {
		clashModeGlobal = "Global"
	}
	clashModeDirect := t.ClashModeDirect
	if clashModeDirect == "" {
		clashModeDirect = "Direct"
	}

	if !t.DisableClashMode {
		options.DNS.Rules = append(options.DNS.Rules, option.DNSRule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultDNSRule{
				ClashMode: clashModeGlobal,
				Server:    DNSDefaultTag,
			},
		}, option.DNSRule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultDNSRule{
				ClashMode: clashModeDirect,
				Server:    DNSLocalTag,
			},
		})
	}
	options.DNS.Rules = append(options.DNS.Rules, t.PreDNSRules...)
	if len(t.CustomDNSRules) == 0 {
		if !t.DisableTrafficBypass {
			if t.DisableRuleSet || (metadata.Version != nil && metadata.Version.LessThan(semver.ParseVersion("1.8.0-alpha.10"))) {
				options.DNS.Rules = append(options.DNS.Rules, option.DNSRule{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultDNSRule{
						Geosite: []string{"geolocation-cn"},
						Server:  DNSLocalTag,
					},
				})
			} else {
				options.DNS.Rules = append(options.DNS.Rules, option.DNSRule{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultDNSRule{
						RuleSet: []string{"geosite-geolocation-cn"},
						Server:  DNSLocalTag,
					},
				})
			}
			if !t.DisableDNSLeak && (metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.9.0-alpha.1"))) {
				options.DNS.Rules = append(options.DNS.Rules, option.DNSRule{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultDNSRule{
						ClashMode: clashModeRule,
						Server:    DNSDefaultTag,
					},
				}, option.DNSRule{
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
									RuleSet: []string{"geoip-cn"},
								},
							},
						},
						Server: DNSLocalTag,
					},
				})
			}
		}
	} else {
		options.DNS.Rules = append(options.DNS.Rules, t.CustomDNSRules...)
	}
	if t.EnableFakeIP {
		options.DNS.Rules = append(options.DNS.Rules, option.DNSRule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultDNSRule{
				QueryType: []option.DNSQueryType{
					option.DNSQueryType(mDNS.TypeA),
					option.DNSQueryType(mDNS.TypeAAAA),
				},
				Server: DNSFakeIPTag,
			},
		})
	}
	return nil
}
