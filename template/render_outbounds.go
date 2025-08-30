package template

import (
	"bytes"
	"regexp"
	"sort"
	"text/template"

	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	"github.com/sagernet/serenity/option"
	"github.com/sagernet/serenity/subscription"
	C "github.com/sagernet/sing-box/constant"
	boxOption "github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
)

func (t *Template) renderOutbounds(metadata M.Metadata, options *boxOption.Options, outbounds [][]boxOption.Outbound, subscriptions []*subscription.Subscription) error {
	disableRuleAction := t.DisableRuleAction || (metadata.Version != nil && metadata.Version.LessThan(semver.ParseVersion("1.11.0-alpha.7")))
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
			Tag:     directTag,
			Type:    C.TypeDirect,
			Options: common.Ptr(common.PtrValueOrDefault(t.CustomDirect)),
		},
		{
			Tag:     defaultTag,
			Type:    C.TypeSelector,
			Options: common.Ptr(common.PtrValueOrDefault(t.CustomSelector)),
		},
	}
	if disableRuleAction {
		options.Outbounds = append(options.Outbounds,
			boxOption.Outbound{
				Tag:     blockTag,
				Type:    C.TypeBlock,
				Options: &boxOption.StubOptions{},
			},
			boxOption.Outbound{
				Tag:     DNSTag,
				Type:    C.TypeDNS,
				Options: &boxOption.StubOptions{},
			},
		)
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
			selectorOptions := common.PtrValueOrDefault(it.CustomSelector)
			selectorOutbound := boxOption.Outbound{
				Type:    C.TypeSelector,
				Tag:     it.Name,
				Options: &selectorOptions,
			}
			selectorOptions.Outbounds = append(selectorOptions.Outbounds, joinOutbounds...)
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
			urltestOptions := common.PtrValueOrDefault(t.CustomURLTest)
			urltestOutbound := boxOption.Outbound{
				Type:    C.TypeURLTest,
				Tag:     urltestTag,
				Options: &urltestOptions,
			}
			urltestOptions.Outbounds = append(urltestOptions.Outbounds, joinOutbounds...)
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
		if extraGroup.Target != option.ExtraGroupTargetSubscription {
			continue
		}
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
			subscriptionTags := common.Filter(common.Map(it.Servers, outboundToString), func(outboundTag string) bool {
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
			})
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
				Tag:  tagPerSubscription,
				Type: extraGroup.Type,
			}
			switch extraGroup.Type {
			case C.TypeSelector:
				selectorOptions := common.PtrValueOrDefault(extraGroup.CustomSelector)
				groupOutboundPerSubscription.Options = &selectorOptions
				selectorOptions.Outbounds = common.Uniq(append(selectorOptions.Outbounds, subscriptionTags...))
				if len(selectorOptions.Outbounds) == 0 {
					continue
				}
			case C.TypeURLTest:
				urltestOptions := common.PtrValueOrDefault(extraGroup.CustomURLTest)
				groupOutboundPerSubscription.Options = &urltestOptions
				urltestOptions.Outbounds = common.Uniq(append(urltestOptions.Outbounds, subscriptionTags...))
				if len(urltestOptions.Outbounds) == 0 {
					continue
				}
			}
			subscriptionGroups[it.Name] = append(subscriptionGroups[it.Name], groupOutboundPerSubscription)
		}
	}
	for _, extraGroup := range t.groups {
		if extraGroup.Target == option.ExtraGroupTargetSubscription {
			continue
		}
		extraTags := groupTags
		for _, group := range subscriptionGroups {
			extraTags = append(extraTags, common.Map(group, outboundToString)...)
		}
		sort.Strings(extraTags)
		if len(extraTags) == 0 || extraGroup.filter != nil || extraGroup.exclude != nil {
			extraTags = append(extraTags, common.Filter(common.FlatMap(subscriptions, func(it *subscription.Subscription) []string {
				return common.Map(it.Servers, outboundToString)
			}), func(outboundTag string) bool {
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
			})...)
		}
		groupOutbound := boxOption.Outbound{
			Tag:  extraGroup.Tag,
			Type: extraGroup.Type,
		}
		switch extraGroup.Type {
		case C.TypeSelector:
			selectorOptions := common.PtrValueOrDefault(extraGroup.CustomSelector)
			groupOutbound.Options = &selectorOptions
			selectorOptions.Outbounds = common.Uniq(append(selectorOptions.Outbounds, extraTags...))
			if len(selectorOptions.Outbounds) == 0 {
				continue
			}
		case C.TypeURLTest:
			urltestOptions := common.PtrValueOrDefault(extraGroup.CustomURLTest)
			groupOutbound.Options = &urltestOptions
			urltestOptions.Outbounds = common.Uniq(append(urltestOptions.Outbounds, extraTags...))
			if len(urltestOptions.Outbounds) == 0 {
				continue
			}
		}
		if extraGroup.Target == option.ExtraGroupTargetDefault {
			defaultGroups = append(defaultGroups, groupOutbound)
		} else {
			globalGroups = append(globalGroups, groupOutbound)
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
	switch outboundOptions := groupOutbound.Options.(type) {
	case *boxOption.SelectorOutboundOptions:
		outboundPtr = &outboundOptions.Outbounds
	case *boxOption.URLTestOutboundOptions:
		outboundPtr = &outboundOptions.Outbounds
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
