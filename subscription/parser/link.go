package parser

import (
	"strings"

	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func ParseSubscriptionLink(link string) (option.Outbound, error) {
	schemeIndex := strings.Index(link, "://")
	if schemeIndex == -1 {
		return option.Outbound{}, E.New("not a link")
	}
	scheme := link[:schemeIndex]
	switch scheme {
	case "ss":
		return ParseShadowsocksLink(link)
	default:
		return option.Outbound{}, E.New("unsupported scheme: ", scheme)
	}
}
