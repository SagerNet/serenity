package template

import (
	"bytes"
	"regexp"
	"text/template"

	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/option"
	"github.com/sagernet/serenity/subscription"
	C "github.com/sagernet/sing-box/constant"
	boxOption "github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
)

func (t *Template) renderOutbounds(metadata M.Metadata, options *boxOption.Options, outbounds [][]boxOption.Outbound, subscriptions []*subscription.Subscription) error {
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
	options.Outbounds = []boxOption.Outbound{
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
	outboundToString := func(it boxOption.Outbound) string {
		return it.Tag
	}
	var globalOutboundTags []string
	if len(outbounds) > 0 {
		for _, outbound := range outbounds {
			options.Outbounds = append(options.Outbounds, outbound...)
		}
		globalOutboundTags = common.Map(outbounds, func(it []boxOption.Outbound) string {
			return it[0].Tag
		})
	}

	var (
		allGroups         []boxOption.Outbound
		allGroupOutbounds []boxOption.Outbound
		groupTags         []string
	)

	for _, it := range subscriptions {
		if len(it.Servers) == 0 {
			continue
		}
		joinOutbounds := common.Map(it.Servers, func(it boxOption.Outbound) string {
			return it.Tag
		})
		if it.GenerateSelector {
			selectorOutbound := boxOption.Outbound{
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
			urltestOutbound := boxOption.Outbound{
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

	var (
		defaultGroups      []boxOption.Outbound
		globalGroups       []boxOption.Outbound
		subscriptionGroups = make(map[string][]boxOption.Outbound)
	)
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
		if extraGroup.Target != option.ExtraGroupTargetSubscription {
			extraTags := common.Filter(common.FlatMap(subscriptions, func(it *subscription.Subscription) []string {
				return common.Map(it.Servers, outboundToString)
			}), myFilter)
			if len(extraTags) == 0 {
				continue
			}
			groupOutbound := boxOption.Outbound{
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
			if extraGroup.Target == option.ExtraGroupTargetDefault {
				defaultGroups = append(defaultGroups, groupOutbound)
			} else {
				globalGroups = append(globalGroups, groupOutbound)
			}
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
				groupOutboundPerSubscription := boxOption.Outbound{
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
				subscriptionGroups[it.Name] = append(subscriptionGroups[it.Name], groupOutboundPerSubscription)
			}
		}
	}

	options.Outbounds = append(options.Outbounds, allGroups...)
	if len(defaultGroups) > 0 {
		options.Outbounds = append(options.Outbounds, defaultGroups...)
	}
	if len(globalGroups) > 0 {
		options.Outbounds = append(options.Outbounds, globalGroups...)
		options.Outbounds = groupJoin(options.Outbounds, defaultTag, false, common.Map(globalGroups, outboundToString)...)
	}
	for _, it := range subscriptions {
		extraGroupOutboundsForSubscription := subscriptionGroups[it.Name]
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

func groupJoin(outbounds []boxOption.Outbound, groupTag string, appendFront bool, groupOutbounds ...string) []boxOption.Outbound {
	groupIndex := common.Index(outbounds, func(it boxOption.Outbound) bool {
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
