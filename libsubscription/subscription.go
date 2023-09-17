package libsubscription

import (
	E "github.com/sagernet/sing/common/exceptions"
)

var subscriptionParsers = []func(string) ([]Server, error){
	ParseBoxSubscription,
	ParseClashSubscription,
	ParseSIP008Subscription,
	ParseRawSubscription,
}

func ParseSubscription(content string) ([]Server, error) {
	var pErr error
	for _, parser := range subscriptionParsers {
		servers, err := parser(content)
		if len(servers) > 0 {
			return servers, nil
		}
		pErr = E.Errors(pErr, err)
	}
	return nil, E.Cause(pErr, "no servers found")
}
