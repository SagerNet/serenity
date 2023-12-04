package filter

import (
	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func init() {
	filters = append(filters, filter180)
}

func filter180(metadata metadata.Metadata, options *option.Options) {
	if metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.8.0-alpha.10")) {
		return
	}
	for index, outbound := range options.Outbounds {
		switch outbound.Type {
		case C.TypeURLTest:
			options.Outbounds[index].URLTestOptions = filter180a10URLTest(outbound.URLTestOptions)
		}
	}
	if metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.8.0-alpha.5")) {
		return
	}
	options.Route.RuleSet = nil
	if options.Route != nil {
		options.Route.Rules = common.Filter(options.Route.Rules, filter180Rule)
	}
	if options.DNS != nil {
		options.DNS.Rules = common.Filter(options.DNS.Rules, filter180DNSRule)
	}
	if metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.8.0-alpha.1")) {
		return
	}
	if options.Route != nil {
		options.Route.Rules = common.Filter(options.Route.Rules, filter180a5Rule)
	}
}

func filter180Rule(it option.Rule) bool {
	return !hasRule([]option.Rule{it}, isRuleSetRule)
}

func filter180DNSRule(it option.DNSRule) bool {
	return !hasDNSRule([]option.DNSRule{it}, isRuleSetDNSRule)
}

func filter180a5Rule(it option.Rule) bool {
	return !hasRule([]option.Rule{it}, isIPIsPrivateRule)
}

func filter180a10URLTest(options option.URLTestOutboundOptions) option.URLTestOutboundOptions {
	options.IdleTimeout = 0
	return options
}
