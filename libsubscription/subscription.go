package libsubscription

import (
	"github.com/sagernet/sing-box/log"
	E "github.com/sagernet/sing/common/exceptions"
)

var subscriptionParsers = []func(string) ([]Server, error){
	ParseClashSubscription,
	ParseSIP008Subscription,
	ParseRawSubscription,
}

func ParseSubscription(content string) ([]Server, error) {
	for _, parser := range subscriptionParsers {
		servers, err := parser(content)
		if len(servers) > 0 {
			return servers, nil
		}
		log.Trace("parse subscription failed: ", err)
	}
	return nil, E.New("no servers found")
}
