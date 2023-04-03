package libsubscription

import (
	"strings"

	E "github.com/sagernet/sing/common/exceptions"
)

func ParseSubscriptionLink(link string) (Server, error) {
	schemeIndex := strings.Index(link, "://")
	if schemeIndex == -1 {
		return Server{}, E.New("not a link")
	}
	scheme := link[:schemeIndex]
	switch scheme {
	case "ss":
		return ParseShadowsocksLink(link)
	default:
		return Server{}, E.New("unsupported scheme: ", scheme)
	}
}
