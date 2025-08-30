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
			Rules:   t.StartRules,
			RuleSet: t.renderRuleSet(t.CustomRuleSet),
		}
	}
	if !t.DisableTrafficBypass {
		t.renderGeoResources(metadata, options)
	}
	disableRuleAction := t.DisableRuleAction || (metadata.Version != nil && metadata.Version.LessThan(semver.ParseVersion("1.11.0-alpha.7")))
	if disableRuleAction {
		options.Route.Rules = append(options.Route.Rules, option.Rule{
			Type: C.RuleTypeLogical,
			LogicalOptions: option.LogicalRule{
				RawLogicalRule: option.RawLogicalRule{
					Mode: C.LogicalTypeOr,
					Rules: []option.Rule{
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultRule{
								RawDefaultRule: option.RawDefaultRule{
									Port: []uint16{53},
								},
							},
						},
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultRule{
								RawDefaultRule: option.RawDefaultRule{
									Protocol: []string{C.ProtocolDNS},
								},
							},
						},
					},
				},
				RuleAction: option.RuleAction{
					Action: C.RuleActionTypeRoute,
					RouteOptions: option.RouteActionOptions{
						Outbound: DNSTag,
					},
				},
			},
		})
	} else {
		options.Route.Rules = append(options.Route.Rules, option.Rule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RuleAction: option.RuleAction{
					Action: C.RuleActionTypeSniff,
				},
			},
		},
			option.Rule{
				Type: C.RuleTypeLogical,
				LogicalOptions: option.LogicalRule{
					RawLogicalRule: option.RawLogicalRule{
						Mode: C.LogicalTypeOr,
						Rules: []option.Rule{
							{
								Type: C.RuleTypeDefault,
								DefaultOptions: option.DefaultRule{
									RawDefaultRule: option.RawDefaultRule{
										Port: []uint16{53},
									},
								},
							},
							{
								Type: C.RuleTypeDefault,
								DefaultOptions: option.DefaultRule{
									RawDefaultRule: option.RawDefaultRule{
										Protocol: []string{C.ProtocolDNS},
									},
								},
							},
						},
					},
					RuleAction: option.RuleAction{
						Action: C.RuleActionTypeHijackDNS,
					},
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
	options.Route.Rules = append(options.Route.Rules, option.Rule{
		Type: C.RuleTypeDefault,
		DefaultOptions: option.DefaultRule{
			RawDefaultRule: option.RawDefaultRule{
				IPIsPrivate: true,
			},
			RuleAction: option.RuleAction{
				Action: C.RuleActionTypeRoute,
				RouteOptions: option.RouteActionOptions{
					Outbound: directTag,
				},
			},
		},
	})
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
				RawDefaultRule: option.RawDefaultRule{
					ClashMode: modeGlobal,
				},
				RuleAction: option.RuleAction{
					Action: C.RuleActionTypeRoute,
					RouteOptions: option.RouteActionOptions{
						Outbound: defaultTag,
					},
				},
			},
		}, option.Rule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RawDefaultRule: option.RawDefaultRule{
					ClashMode: modeDirect,
				},
				RuleAction: option.RuleAction{
					Action: C.RuleActionTypeRoute,
					RouteOptions: option.RouteActionOptions{
						Outbound: directTag,
					},
				},
			},
		})
	}
	if !disableRuleAction {
		options.Route.Rules = append(options.Route.Rules, option.Rule{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RuleAction: option.RuleAction{
					Action: C.RuleActionTypeResolve,
				},
			},
		})
	}
	options.Route.Rules = append(options.Route.Rules, t.PreRules...)
	if len(t.CustomRules) == 0 {
		if !t.DisableTrafficBypass {
			options.Route.Rules = append(options.Route.Rules, option.Rule{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					RawDefaultRule: option.RawDefaultRule{
						RuleSet: []string{"geosite-geolocation-cn"},
					},
					RuleAction: option.RuleAction{
						Action: C.RuleActionTypeRoute,
						RouteOptions: option.RouteActionOptions{
							Outbound: directTag,
						},
					},
				},
			}, option.Rule{
				Type: C.RuleTypeLogical,
				LogicalOptions: option.LogicalRule{
					RawLogicalRule: option.RawLogicalRule{
						Mode: C.LogicalTypeAnd,
						Rules: []option.Rule{
							{
								Type: C.RuleTypeDefault,
								DefaultOptions: option.DefaultRule{
									RawDefaultRule: option.RawDefaultRule{
										RuleSet: []string{"geoip-cn"},
									},
								},
							},
							{
								Type: C.RuleTypeDefault,
								DefaultOptions: option.DefaultRule{
									RawDefaultRule: option.RawDefaultRule{
										RuleSet: []string{"geosite-geolocation-!cn"},
										Invert:  true,
									},
								},
							},
						},
					},
					RuleAction: option.RuleAction{
						Action: C.RuleActionTypeRoute,
						RouteOptions: option.RouteActionOptions{
							Outbound: directTag,
						},
					},
				},
			})
		}
	} else {
		options.Route.Rules = append(options.Route.Rules, t.CustomRules...)
	}
	if !t.DisableTrafficBypass && !t.DisableDefaultRules {
		blockTag := t.BlockTag
		if blockTag == "" {
			blockTag = DefaultBlockTag
		}
		options.Route.Rules = append(options.Route.Rules, option.Rule{
			Type: C.RuleTypeLogical,
			LogicalOptions: option.LogicalRule{
				RawLogicalRule: option.RawLogicalRule{
					Mode: C.LogicalTypeOr,
					Rules: []option.Rule{
						{
							Type: C.RuleTypeDefault,
							DefaultOptions: option.DefaultRule{
								RawDefaultRule: option.RawDefaultRule{
									Network: []string{N.NetworkUDP},
									Port:    []uint16{443},
								},
							},
						},
					},
				},
				RuleAction: option.RuleAction{
					Action: C.RuleActionTypeRoute,
					RouteOptions: option.RouteActionOptions{
						Outbound: blockTag,
					},
				},
			},
		})
	}
	if metadata.Version != nil && metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.12.0-alpha.1")) {
		options.Route.DefaultDomainResolver = &option.DomainResolveOptions{
			Server: DNSLocalTag,
		}
	}
	return nil
}
