package filter

import (
	"context"
	"net/netip"

	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
)

func init() {
	filters = append(filters, filter1100)
}

func filter1100(metadata metadata.Metadata, options *option.Options) error {
	if metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.10.0-alpha.19")) {
		return nil
	}
	var newRuleSets []option.RuleSet
	var inlineRuleSets []option.RuleSet
	for _, ruleSet := range options.Route.RuleSet {
		if ruleSet.Type == C.RuleSetTypeInline {
			inlineRuleSets = append(inlineRuleSets, ruleSet)
		} else {
			newRuleSets = append(newRuleSets, ruleSet)
		}
	}
	options.Route.RuleSet = newRuleSets
	if len(inlineRuleSets) > 0 {
		var (
			currentRules []option.Rule
			newRules     []option.Rule
		)
		currentRules = options.Route.Rules
		for _, inlineRuleSet := range inlineRuleSets {
			for i, rule := range currentRules {
				newRuleItems, err := expandInlineRule(inlineRuleSet, rule)
				if err != nil {
					return E.Cause(err, "expand rule[", i, "]")
				}
				newRules = append(newRules, newRuleItems...)
			}
			currentRules = newRules
			newRules = newRules[:0]
		}
		options.Route.Rules = currentRules

		var (
			currentDNSRules []option.DNSRule
			newDNSRules     []option.DNSRule
		)
		currentDNSRules = options.DNS.Rules
		for _, inlineRuleSet := range inlineRuleSets {
			for i, rule := range currentDNSRules {
				newRuleItems, err := expandInlineDNSRule(inlineRuleSet, rule)
				if err != nil {
					return E.Cause(err, "expand dns rule[", i, "]")
				}
				newDNSRules = append(newDNSRules, newRuleItems...)
			}
			currentDNSRules = newDNSRules
			newDNSRules = newDNSRules[:0]
		}
		options.DNS.Rules = currentDNSRules
	}
	options.Route.Rules = common.Filter(options.Route.Rules, filter1100Rule)
	options.DNS.Rules = common.Filter(options.DNS.Rules, filter1100DNSRule)
	if metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.10.0-alpha.13")) {
		return nil
	}
	if len(options.Inbounds) > 0 {
		newInbounds := make([]option.Inbound, 0, len(options.Inbounds))
		for _, inbound := range options.Inbounds {
			if inbound.Type == C.TypeTun {
				tunOptions := inbound.Options.(*option.TunInboundOptions)
				tunOptions.AutoRedirect = false
				tunOptions.RouteAddressSet = nil
				tunOptions.RouteExcludeAddressSet = nil
				//nolint:staticcheck
				//goland:noinspection GoDeprecation
				if len(tunOptions.Address) > 0 {
					tunOptions.Inet4Address = append(tunOptions.Inet4Address, common.Filter(tunOptions.Address, func(it netip.Prefix) bool {
						return it.Addr().Is4()
					})...)
					tunOptions.Inet6Address = append(tunOptions.Inet6Address, common.Filter(tunOptions.Address, func(it netip.Prefix) bool {
						return it.Addr().Is6()
					})...)
					tunOptions.Address = nil
				}
				//nolint:staticcheck
				//goland:noinspection GoDeprecation
				if len(tunOptions.RouteAddress) > 0 {
					tunOptions.Inet4RouteAddress = append(tunOptions.Inet4RouteAddress, common.Filter(tunOptions.RouteAddress, func(it netip.Prefix) bool {
						return it.Addr().Is4()
					})...)
					tunOptions.Inet6RouteAddress = append(tunOptions.Inet6RouteAddress, common.Filter(tunOptions.RouteAddress, func(it netip.Prefix) bool {
						return it.Addr().Is6()
					})...)
					tunOptions.RouteAddress = nil
				}
				//nolint:staticcheck
				//goland:noinspection GoDeprecation
				if len(tunOptions.RouteExcludeAddress) > 0 {
					tunOptions.Inet4RouteExcludeAddress = append(tunOptions.Inet4RouteExcludeAddress, common.Filter(tunOptions.RouteExcludeAddress, func(it netip.Prefix) bool {
						return it.Addr().Is4()
					})...)
					tunOptions.Inet6RouteExcludeAddress = append(tunOptions.Inet6RouteExcludeAddress, common.Filter(tunOptions.RouteExcludeAddress, func(it netip.Prefix) bool {
						return it.Addr().Is6()
					})...)
					tunOptions.RouteExcludeAddress = nil
				}
			}
			newInbounds = append(newInbounds, inbound)
		}
		options.Inbounds = newInbounds
	}
	return nil
}

func expandInlineRule(ruleSet option.RuleSet, rule option.Rule) ([]option.Rule, error) {
	var (
		newRules []option.Rule
		err      error
	)
	if rule.Type == C.RuleTypeLogical {
		for i := range rule.LogicalOptions.Rules {
			newRules, err = expandInlineRule(ruleSet, rule.LogicalOptions.Rules[i])
			if err != nil {
				return nil, E.Cause(err, "[", i, "]")
			}
			newRules = append(newRules, newRules...)
		}
		rule.LogicalOptions.Rules = newRules
		return []option.Rule{rule}, nil
	}
	if !common.Contains(rule.DefaultOptions.RuleSet, ruleSet.Tag) {
		return []option.Rule{rule}, nil
	}
	rule.DefaultOptions.RuleSet = common.Filter(rule.DefaultOptions.RuleSet, func(it string) bool {
		return it != ruleSet.Tag
	})
	for i, hRule := range ruleSet.InlineOptions.Rules {
		var (
			rawRule json.RawMessage
			newRule option.Rule
		)
		rawRule, err = json.Marshal(hRule)
		if err != nil {
			return nil, E.Cause(err, "marshal inline rule ", ruleSet.Tag, "[", i, "]")
		}
		newRule, err = badjson.MergeFromSource(context.Background(), rawRule, rule, false)
		if err != nil {
			return nil, E.Cause(err, "merge inline rule ", ruleSet.Tag, "[", i, "]")
		}
		newRules = append(newRules, newRule)
	}
	return newRules, nil
}

func expandInlineDNSRule(ruleSet option.RuleSet, rule option.DNSRule) ([]option.DNSRule, error) {
	var (
		newRules []option.DNSRule
		err      error
	)
	if rule.Type == C.RuleTypeLogical {
		for i := range rule.LogicalOptions.Rules {
			newRules, err = expandInlineDNSRule(ruleSet, rule.LogicalOptions.Rules[i])
			if err != nil {
				return nil, E.Cause(err, "[", i, "]")
			}
			newRules = append(newRules, newRules...)
		}
		rule.LogicalOptions.Rules = newRules
		return []option.DNSRule{rule}, nil
	}
	if !common.Contains(rule.DefaultOptions.RuleSet, ruleSet.Tag) {
		return []option.DNSRule{rule}, nil
	}
	rule.DefaultOptions.RuleSet = common.Filter(rule.DefaultOptions.RuleSet, func(it string) bool {
		return it != ruleSet.Tag
	})
	for i, hRule := range ruleSet.InlineOptions.Rules {
		var (
			rawRule json.RawMessage
			newRule option.DNSRule
		)
		rawRule, err = json.Marshal(hRule)
		if err != nil {
			return nil, E.Cause(err, "marshal inline rule ", ruleSet.Tag, "[", i, "]")
		}
		newRule, err = badjson.MergeFromSource(context.Background(), rawRule, rule, false)
		if err != nil {
			return nil, E.Cause(err, "merge inline rule ", ruleSet.Tag, "[", i, "]")
		}
		newRules = append(newRules, newRule)
	}
	return newRules, nil
}

func filter1100Rule(it option.Rule) bool {
	return !hasRule([]option.Rule{it}, func(it option.DefaultRule) bool {
		return it.RuleSetIPCIDRMatchSource
	})
}

func filter1100DNSRule(it option.DNSRule) bool {
	return !hasDNSRule([]option.DNSRule{it}, func(it option.DefaultDNSRule) bool {
		return it.RuleSetIPCIDRMatchSource || it.RuleSetIPCIDRAcceptEmpty
	})
}
