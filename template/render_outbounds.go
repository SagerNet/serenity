package template

import (
	"bytes"
	"regexp"
	"text/template"

	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/subscription"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
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
	outboundToString := func(it option.Outbound) string {
		return it.Tag
	}
	var (
		globalOutbounds    []option.Outbound
		globalOutboundTags []string
	)
	if len(outbounds) > 0 {
		for _, outbound := range outbounds {
			options.Outbounds = append(options.Outbounds, outbound...)
		}
		globalOutbounds = common.Map(outbounds, func(it []option.Outbound) option.Outbound {
			return it[0]
		})
		globalOutboundTags = common.Map(globalOutbounds, outboundToString)
	}

	var (
		allGroups         []option.Outbound
		allGroupOutbounds []option.Outbound
		groupTags         []string
	)

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
			groupTags = append(groupTags, selectorOutbound.Tag)
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
			groupTags = append(groupTags, urltestOutbound.Tag)
		}
		if !it.GenerateSelector && !it.GenerateURLTest {
			globalOutboundTags = append(globalOutboundTags, joinOutbounds...)
		}
		allGroupOutbounds = append(allGroupOutbounds, it.Servers...)
	}

	globalOutbounds = append(globalOutbounds, allGroups...)
	globalOutbounds = append(globalOutbounds, allGroupOutbounds...)

	allExtraGroups := make(map[string][]option.Outbound)
	for _, extraGroup := range t.groups {
		myFilter := func(outboundTag string) bool {
			if len(extraGroup.filter) > 0 {
				if !common.Any(extraGroup.filter, func(it *regexp.Regexp) bool {
					return it.MatchString(outboundTag)
				}) {
					return false
				}
			}
			if len(extraGroup.exclude) > 0 {
				if common.Any(extraGroup.exclude, func(it *regexp.Regexp) bool {
					return it.MatchString(outboundTag)
				}) {
					return false
				}
			}
			return true
		}
		if !extraGroup.PerSubscription {
			var extraTags []string
			if extraGroup.ExcludeOutbounds {
				extraTags = common.Filter(common.FlatMap(subscriptions, func(it *subscription.Subscription) []string {
					return common.Map(it.Servers, outboundToString)
				}), myFilter)
			} else {
				extraTags = common.Filter(common.Map(globalOutbounds, outboundToString), myFilter)
			}
			if len(extraTags) == 0 {
				continue
			}
			groupOutbound := option.Outbound{
				Tag:             extraGroup.Tag,
				Type:            extraGroup.Type,
				SelectorOptions: common.PtrValueOrDefault(extraGroup.CustomSelector),
				URLTestOptions:  common.PtrValueOrDefault(extraGroup.CustomURLTest),
			}
			switch extraGroup.Type {
			case C.TypeSelector:
				groupOutbound.SelectorOptions.Outbounds = append(groupOutbound.SelectorOptions.Outbounds, extraTags...)
			case C.TypeURLTest:
				groupOutbound.URLTestOptions.Outbounds = append(groupOutbound.URLTestOptions.Outbounds, extraTags...)
			}
			allExtraGroups[""] = append(allExtraGroups[""], groupOutbound)
		} else {
			tmpl := template.New("tag")
			if extraGroup.TagPerSubscription != "" {
				_, err := tmpl.Parse(extraGroup.TagPerSubscription)
				if err != nil {
					return E.Cause(err, "parse `tag_per_subscription`: ", extraGroup.TagPerSubscription)
				}
			} else {
				common.Must1(tmpl.Parse("{{ .tag }} ({{ .subscription_name }})"))
			}
			var outboundTags []string
			if !extraGroup.ExcludeOutbounds {
				outboundTags = common.Filter(common.FlatMap(outbounds, func(it []option.Outbound) []string {
					return common.Map(it, outboundToString)
				}), myFilter)
			}
			if len(outboundTags) > 0 {
				groupOutbound := option.Outbound{
					Tag:             extraGroup.Tag,
					Type:            extraGroup.Type,
					SelectorOptions: common.PtrValueOrDefault(extraGroup.CustomSelector),
					URLTestOptions:  common.PtrValueOrDefault(extraGroup.CustomURLTest),
				}
				switch extraGroup.Type {
				case C.TypeSelector:
					groupOutbound.SelectorOptions.Outbounds = append(groupOutbound.SelectorOptions.Outbounds, outboundTags...)
				case C.TypeURLTest:
					groupOutbound.URLTestOptions.Outbounds = append(groupOutbound.URLTestOptions.Outbounds, outboundTags...)
				}
				allExtraGroups[""] = append(allExtraGroups[""], groupOutbound)
			}
			for _, it := range subscriptions {
				subscriptionTags := common.Filter(common.Map(it.Servers, outboundToString), myFilter)
				if len(subscriptionTags) == 0 {
					continue
				}
				var tagPerSubscription string
				if len(outboundTags) == 0 && len(subscriptions) == 1 {
					tagPerSubscription = extraGroup.Tag
				} else {
					var buffer bytes.Buffer
					err := tmpl.Execute(&buffer, map[string]interface{}{
						"tag":               extraGroup.Tag,
						"subscription_name": it.Name,
					})
					if err != nil {
						return E.Cause(err, "generate tag for extra group: tag=", extraGroup.Tag, ", subscription=", it.Name)
					}
					tagPerSubscription = buffer.String()
				}
				groupOutboundPerSubscription := option.Outbound{
					Tag:             tagPerSubscription,
					Type:            extraGroup.Type,
					SelectorOptions: common.PtrValueOrDefault(extraGroup.CustomSelector),
					URLTestOptions:  common.PtrValueOrDefault(extraGroup.CustomURLTest),
				}
				switch extraGroup.Type {
				case C.TypeSelector:
					groupOutboundPerSubscription.SelectorOptions.Outbounds = append(groupOutboundPerSubscription.SelectorOptions.Outbounds, subscriptionTags...)
				case C.TypeURLTest:
					groupOutboundPerSubscription.URLTestOptions.Outbounds = append(groupOutboundPerSubscription.URLTestOptions.Outbounds, subscriptionTags...)
				}
				allExtraGroups[it.Name] = append(allExtraGroups[it.Name], groupOutboundPerSubscription)
			}
		}
	}

	options.Outbounds = append(options.Outbounds, allGroups...)

	defaultExtraGroupOutbounds := allExtraGroups[""]
	if len(defaultExtraGroupOutbounds) > 0 {
		options.Outbounds = append(options.Outbounds, defaultExtraGroupOutbounds...)
		options.Outbounds = groupJoin(options.Outbounds, defaultTag, false, common.Map(defaultExtraGroupOutbounds, outboundToString)...)
	}
	for _, it := range subscriptions {
		extraGroupOutboundsForSubscription := allExtraGroups[it.Name]
		if len(extraGroupOutboundsForSubscription) > 0 {
			options.Outbounds = append(options.Outbounds, extraGroupOutboundsForSubscription...)
			options.Outbounds = groupJoin(options.Outbounds, it.Name, true, common.Map(extraGroupOutboundsForSubscription, outboundToString)...)
		}
	}
	options.Outbounds = groupJoin(options.Outbounds, defaultTag, false, groupTags...)
	options.Outbounds = groupJoin(options.Outbounds, defaultTag, false, globalOutboundTags...)

	options.Outbounds = append(options.Outbounds, allGroupOutbounds...)
	return nil
}

func groupJoin(outbounds []option.Outbound, groupTag string, appendFront bool, groupOutbounds ...string) []option.Outbound {
	groupIndex := common.Index(outbounds, func(it option.Outbound) bool {
		return it.Tag == groupTag
	})
	if groupIndex == -1 {
		return outbounds
	}
	groupOutbound := outbounds[groupIndex]
	var outboundPtr *[]string
	switch groupOutbound.Type {
	case C.TypeSelector:
		outboundPtr = &groupOutbound.SelectorOptions.Outbounds
	case C.TypeURLTest:
		outboundPtr = &groupOutbound.URLTestOptions.Outbounds
	default:
		panic(F.ToString("unexpected group type: ", groupOutbound.Type))
	}
	if appendFront {
		*outboundPtr = append(groupOutbounds, *outboundPtr...)
	} else {
		*outboundPtr = append(*outboundPtr, groupOutbounds...)
	}
	*outboundPtr = common.Dup(*outboundPtr)
	outbounds[groupIndex] = groupOutbound
	return outbounds
}
