package template

import (
	"regexp"

	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/subscription"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func (t *Template) renderOutbounds(metadata M.Metadata, options *option.Options, outbounds [][]option.Outbound, subscriptions []*subscription.Subscription) error {
	defaultTag := t.DefaultTag
	if defaultTag == "" {
		defaultTag = DefaultDefaultTag
	}
	options.Route.Final = defaultTag
	directTag := t.DirectTag
	if directTag == "" {
		directTag = DefaultDirectTag
	}
	blockTag := t.BlockTag
	if blockTag == "" {
		blockTag = DefaultBlockTag
	}
	options.Outbounds = []option.Outbound{
		{
			Tag:           directTag,
			Type:          C.TypeDirect,
			DirectOptions: common.PtrValueOrDefault(t.CustomDirect),
		},
		{
			Tag:  blockTag,
			Type: C.TypeBlock,
		},
		{
			Tag:  DNSTag,
			Type: C.TypeDNS,
		},
		{
			Tag:             defaultTag,
			Type:            C.TypeSelector,
			SelectorOptions: common.PtrValueOrDefault(t.CustomSelector),
		},
	}
	urlTestTag := t.URLTestTag
	if urlTestTag == "" {
		urlTestTag = DefaultURLTestTag
	}
	if t.GenerateGlobalURLTest {
		options.Outbounds = append(options.Outbounds, option.Outbound{
			Tag:            urlTestTag,
			Type:           C.TypeURLTest,
			URLTestOptions: common.PtrValueOrDefault(t.CustomURLTest),
		})
	}
	globalJoin := func(groupOutbounds ...string) {
		options.Outbounds = groupJoin(options.Outbounds, defaultTag, groupOutbounds...)
		if t.GenerateGlobalURLTest {
			options.Outbounds = groupJoin(options.Outbounds, urlTestTag, groupOutbounds...)
		}
	}

	var globalOutbounds []option.Outbound
	if len(outbounds) > 0 {
		for _, outbound := range outbounds {
			options.Outbounds = append(options.Outbounds, outbound...)
		}
		globalOutbounds = common.Map(outbounds, func(it []option.Outbound) option.Outbound {
			return it[0]
		})
		globalJoin(common.Map(globalOutbounds, func(it option.Outbound) string {
			return it.Tag
		})...)
	}

	var allGroups []option.Outbound
	var allGroupOutbounds []option.Outbound

	for _, it := range subscriptions {
		if len(it.Servers) == 0 {
			continue
		}
		joinOutbounds := common.Map(it.Servers, func(it option.Outbound) string {
			return it.Tag
		})
		if it.GenerateSelector {
			selectorOutbound := option.Outbound{
				Type:            C.TypeSelector,
				Tag:             it.Name,
				SelectorOptions: common.PtrValueOrDefault(it.CustomSelector),
			}
			selectorOutbound.SelectorOptions.Outbounds = append(selectorOutbound.SelectorOptions.Outbounds, joinOutbounds...)
			allGroups = append(allGroups, selectorOutbound)
			globalJoin(it.Name)
		}
		if it.GenerateURLTest {
			var urltestTag string
			if !it.GenerateSelector {
				urltestTag = it.Name
			} else if it.URLTestTagSuffix != "" {
				urltestTag = it.Name + " " + it.URLTestTagSuffix
			} else {
				urltestTag = it.Name + " - URLTest"
			}
			urltestOutbound := option.Outbound{
				Type:           C.TypeURLTest,
				Tag:            urltestTag,
				URLTestOptions: common.PtrValueOrDefault(t.CustomURLTest),
			}
			urltestOutbound.URLTestOptions.Outbounds = append(urltestOutbound.URLTestOptions.Outbounds, joinOutbounds...)
			allGroups = append(allGroups, urltestOutbound)
			globalJoin(urltestTag)
		}
		if !it.GenerateSelector && !it.GenerateURLTest {
			globalJoin(joinOutbounds...)
		}
		allGroupOutbounds = append(allGroupOutbounds, it.Servers...)
	}

	globalOutbounds = append(globalOutbounds, allGroups...)
	globalOutbounds = append(globalOutbounds, allGroupOutbounds...)

	for _, group := range t.groups {
		var extraTags []string
		for _, groupOutbound := range globalOutbounds {
			if len(group.filter) > 0 {
				if !common.Any(group.filter, func(it *regexp.Regexp) bool {
					return it.MatchString(groupOutbound.Tag)
				}) {
					continue
				}
			}
			if len(group.exclude) > 0 {
				if common.Any(group.exclude, func(it *regexp.Regexp) bool {
					return it.MatchString(groupOutbound.Tag)
				}) {
					continue
				}
			}
			extraTags = append(extraTags, groupOutbound.Tag)
		}
		if len(extraTags) == 0 {
			continue
		}
		groupOutbound := option.Outbound{
			Tag:             group.Tag,
			Type:            group.Type,
			SelectorOptions: common.PtrValueOrDefault(group.CustomSelector),
			URLTestOptions:  common.PtrValueOrDefault(group.CustomURLTest),
		}
		switch group.Type {
		case C.TypeSelector:
			groupOutbound.SelectorOptions.Outbounds = append(groupOutbound.SelectorOptions.Outbounds, extraTags...)
		case C.TypeURLTest:
			groupOutbound.URLTestOptions.Outbounds = append(groupOutbound.URLTestOptions.Outbounds, extraTags...)
		}
		options.Outbounds = append(options.Outbounds, groupOutbound)
	}

	options.Outbounds = append(options.Outbounds, allGroups...)
	options.Outbounds = append(options.Outbounds, allGroupOutbounds...)

	return nil
}

func groupJoin(outbounds []option.Outbound, groupTag string, groupOutbounds ...string) []option.Outbound {
	groupIndex := common.Index(outbounds, func(it option.Outbound) bool {
		return it.Tag == groupTag
	})
	if groupIndex == -1 {
		return outbounds
	}
	groupOutbound := outbounds[groupIndex]
	switch groupOutbound.Type {
	case C.TypeSelector:
		groupOutbound.SelectorOptions.Outbounds = common.Dup(append(groupOutbound.SelectorOptions.Outbounds, groupOutbounds...))
	case C.TypeURLTest:
		groupOutbound.URLTestOptions.Outbounds = common.Dup(append(groupOutbound.URLTestOptions.Outbounds, groupOutbounds...))
	}
	outbounds[groupIndex] = groupOutbound
	return outbounds
}
