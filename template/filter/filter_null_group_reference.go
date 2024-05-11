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

func filterNullGroupReference(metadata M.Metadata, options *option.Options) {
	outboundTags := common.Map(options.Outbounds, func(it option.Outbound) string {
		return it.Tag
	})
	for i, outbound := range options.Outbounds {
		switch outbound.Type {
		case C.TypeSelector:
			outbound.SelectorOptions.Outbounds = common.Filter(outbound.SelectorOptions.Outbounds, func(outbound string) bool {
				return common.Contains(outboundTags, outbound)
			})
		case C.TypeURLTest:
			outbound.URLTestOptions.Outbounds = common.Filter(outbound.URLTestOptions.Outbounds, func(outbound string) bool {
				return common.Contains(outboundTags, outbound)
			})
		default:
			continue
		}
		options.Outbounds[i] = outbound
	}
}
