package filter

import (
	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func init() {
	filters = append(filters, filter190)
}

func filter190(metadata metadata.Metadata, options *option.Options) {
	if metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.9.0-alpha.1")) {
		return
	}
	if options.DNS == nil || len(options.DNS.Rules) == 0 {
		return
	}
	options.DNS.Rules = common.Filter(options.DNS.Rules, filter190DNSRule)
	if metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.9.0-alpha.10")) {
		return
	}
	for _, inbound := range options.Inbounds {
		switch inbound.Type {
		case C.TypeTun:
			if inbound.TunOptions.Platform == nil || inbound.TunOptions.Platform.HTTPProxy == nil {
				continue
			}
			httpProxy := inbound.TunOptions.Platform.HTTPProxy
			if len(httpProxy.BypassDomain) > 0 || len(httpProxy.MatchDomain) > 0 {
				httpProxy.BypassDomain = nil
				httpProxy.MatchDomain = nil
			}
		}
	}
}

func filter190DNSRule(it option.DNSRule) bool {
	return !hasDNSRule([]option.DNSRule{it}, isAddressFilterRule)
}

func isAddressFilterRule(it option.DefaultDNSRule) bool {
	return len(it.GeoIP) > 0 || len(it.IPCIDR) > 0 || it.IPIsPrivate
}
