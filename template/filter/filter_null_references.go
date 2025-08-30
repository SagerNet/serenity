package filter

import (
	M "github.com/sagernet/serenity/common/metadata"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func init() {
	filters = append(filters, filterNullGroupReference)
}

func filterNullGroupReference(metadata M.Metadata, options *option.Options) error {
	outboundTags := common.Map(options.Outbounds, func(it option.Outbound) string {
		return it.Tag
	})
	for _, outbound := range options.Outbounds {
		switch outboundOptions := outbound.Options.(type) {
		case *option.SelectorOutboundOptions:
			outboundOptions.Outbounds = common.Filter(outboundOptions.Outbounds, func(outbound string) bool {
				return common.Contains(outboundTags, outbound)
			})
		case *option.URLTestOutboundOptions:
			outboundOptions.Outbounds = common.Filter(outboundOptions.Outbounds, func(outbound string) bool {
				return common.Contains(outboundTags, outbound)
			})
		default:
			continue
		}
	}
	options.Route.Rules = common.Filter(options.Route.Rules, func(it option.Rule) bool {
		switch it.Type {
		case C.RuleTypeDefault:
			if it.DefaultOptions.Action != C.RuleActionTypeRoute {
				return true
			}
			return common.Contains(outboundTags, it.DefaultOptions.RouteOptions.Outbound)
		case C.RuleTypeLogical:
			if it.LogicalOptions.Action != C.RuleActionTypeRoute {
				return true
			}
			return common.Contains(outboundTags, it.LogicalOptions.RouteOptions.Outbound)
		default:
			panic("no")
		}
	})
	return nil
}
