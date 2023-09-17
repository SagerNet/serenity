package libsubscription

import (
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

func ParseBoxSubscription(content string) ([]Server, error) {
	var options option.Options
	err := options.UnmarshalJSON([]byte(content))
	if err != nil {
		return nil, err
	}
	if len(options.Outbounds) == 0 {
		return nil, E.New("no servers found")
	}
	options.Outbounds = common.Filter(options.Outbounds, func(it option.Outbound) bool {
		switch it.Type {
		case C.TypeDirect, C.TypeBlock, C.TypeDNS, C.TypeSelector, C.TypeURLTest:
			return false
		default:
			return true
		}
	})
	return common.Map(options.Outbounds, func(it option.Outbound) Server {
		return Server{
			Name:      it.Tag,
			Outbounds: []option.Outbound{it},
		}
	}), nil
}
