package filter

import (
	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func init() {
	filters = append(filters, filter1100)
}

func filter1100(metadata metadata.Metadata, options *option.Options) {
	if metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.10.0-alpha.2")) {
		return
	}
	newInbounds := make([]option.Inbound, 0, len(options.Inbounds))
	for _, inbound := range options.Inbounds {
		if inbound.Type == C.TypeTun && inbound.TunOptions.AutoRedirect {
			inbound.TunOptions.AutoRedirect = false
		}
		newInbounds = append(newInbounds, inbound)
	}
	options.Inbounds = newInbounds
}
