package template

import (
	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	N "github.com/sagernet/sing/common/network"
)

func (t *Template) renderRoute(metadata M.Metadata, options *option.Options) error {
	if options.Route == nil {
		options.Route = &option.RouteOptions{
			GeoIP:   t.CustomGeoIP,
			Geosite: t.CustomGeosite,
			RuleSet: t.CustomRuleSet,
		}
	}
	if !t.DisableTrafficBypass {
		t.renderGeoResources(metadata, options)
	}
	disable18Features := metadata.Version != nil && metadata.Version.LessThan(semver.ParseVersion("1.8.0-alpha.10"))
	options.Route.Rules = []option.Rule{
		{
			Type: C.RuleTypeLogical,
			LogicalOptions: option.LogicalRule{
				Mode: C.LogicalTypeOr,
				Rules: []option.Rule{
					{
						Type: C.RuleTypeDefault,
						DefaultOptions: option.DefaultRule{
							Network: []string{N.NetworkUDP},
							Port:    []uint16{53},
						},
					},
					{
						Type: C.RuleTypeDefault,
						DefaultOptions: option.DefaultRule{
							Protocol: []string{C.ProtocolDNS},
						},
					},
				},
				Outbound: DNSTag,
			},
		},
	}
	if !t.DisableTrafficBypass && !t.DisableDefaultRules {
		options.Route.Rules = append(options.Route.Rules, option.Rule{
			Type: C.RuleTypeLogical,
			LogicalOptions: option.LogicalRule{
				Mode: C.LogicalTypeOr,
				Rules: []option.Rule{
					{
						Type: C.RuleTypeDefault,
						DefaultOptions: option.DefaultRule{
							Network: []string{N.NetworkUDP},
							Port:    []uint16{443},
						},
					},
					{
						Type: C.RuleTypeDefault,
						DefaultOptions: option.DefaultRule{
							Protocol: []string{C.ProtocolSTUN},
						},
					},
				},
				Outbound: BlockTag,
			},
		})
	}
	directTag := t.DirectTag
	defaultTag := t.DefaultTag
	if directTag == "" {
		directTag = DefaultDirectTag
	}
	if defaultTag == "" {
		defaultTag = DefaultDefaultTag
	}
	if disable18Features {
		options.Route.Rules = append(options.Route.Rules, option.Rule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				GeoIP:    []string{"private"},
				Outbound: directTag,
			},
		})
	} else {
		options.Route.Rules = append(options.Route.Rules, option.Rule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				IPIsPrivate: true,
				Outbound:    directTag,
			},
		})
	}
	if !t.DisableClashMode {
		modeGlobal := t.ClashModeGlobal
		modeDirect := t.ClashModeDirect
		if modeGlobal == "" {
			modeGlobal = "Global"
		}
		if modeDirect == "" {
			modeDirect = "Direct"
		}
		options.Route.Rules = append(options.Route.Rules, option.Rule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				ClashMode: modeGlobal,
				Outbound:  defaultTag,
			},
		}, option.Rule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				ClashMode: modeDirect,
				Outbound:  directTag,
			},
		})
	}
	options.Route.Rules = append(options.Route.Rules, t.PreRules...)
	if len(t.CustomRules) == 0 {
		if !t.DisableTrafficBypass {
			if t.DisableRuleSet || disable18Features {
				options.Route.Rules = append(options.Route.Rules, option.Rule{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						GeoIP:    []string{"cn"},
						Geosite:  []string{"geolocation-cn"},
						Outbound: directTag,
					},
				})
			} else {
				options.Route.Rules = append(options.Route.Rules, option.Rule{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						RuleSet:  []string{"geoip-cn", "geosite-geolocation-cn"},
						Outbound: directTag,
					},
				})
			}
		}
	} else {
		options.Route.Rules = append(options.Route.Rules, t.CustomRules...)
	}
	return nil
}
