package parser

import (
	"context"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
)

func ParseBoxSubscription(ctx context.Context, content string) ([]option.Outbound, error) {
	options, err := json.UnmarshalExtendedContext[option.Options](ctx, []byte(content))
	if err != nil {
		return nil, err
	}
	options.Outbounds = common.Filter(options.Outbounds, func(it option.Outbound) bool {
		switch it.Type {
		case C.TypeDirect, C.TypeBlock, C.TypeDNS, C.TypeSelector, C.TypeURLTest:
			return false
		default:
			return true
		}
	})
	if len(options.Outbounds) == 0 {
		return nil, E.New("no servers found")
	}
	return options.Outbounds, nil
}
