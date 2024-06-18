package filter

import (
	"net/netip"

	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func init() {
	filters = append(filters, filter1100)
}

func filter1100(metadata metadata.Metadata, options *option.Options) {
	if metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.10.0-alpha.13")) {
		return
	}
	newInbounds := make([]option.Inbound, 0, len(options.Inbounds))
	for _, inbound := range options.Inbounds {
		if inbound.Type == C.TypeTun {
			inbound.TunOptions.AutoRedirect = false
			inbound.TunOptions.RouteAddressSet = nil
			inbound.TunOptions.RouteExcludeAddressSet = nil
			//nolint:staticcheck
			//goland:noinspection GoDeprecation
			if len(inbound.TunOptions.Address) > 0 {
				inbound.TunOptions.Inet4Address = append(inbound.TunOptions.Inet4Address, common.Filter(inbound.TunOptions.Address, func(it netip.Prefix) bool {
					return it.Addr().Is4()
				})...)
				inbound.TunOptions.Inet6Address = append(inbound.TunOptions.Inet6Address, common.Filter(inbound.TunOptions.Address, func(it netip.Prefix) bool {
					return it.Addr().Is6()
				})...)
			}
			//nolint:staticcheck
			//goland:noinspection GoDeprecation
			if len(inbound.TunOptions.RouteAddress) > 0 {
				inbound.TunOptions.Inet4RouteAddress = append(inbound.TunOptions.Inet4RouteAddress, common.Filter(inbound.TunOptions.RouteAddress, func(it netip.Prefix) bool {
					return it.Addr().Is4()
				})...)
				inbound.TunOptions.Inet6RouteAddress = append(inbound.TunOptions.Inet6RouteAddress, common.Filter(inbound.TunOptions.RouteAddress, func(it netip.Prefix) bool {
					return it.Addr().Is6()
				})...)
			}
			//nolint:staticcheck
			//goland:noinspection GoDeprecation
			if len(inbound.TunOptions.RouteExcludeAddress) > 0 {
				inbound.TunOptions.Inet4RouteExcludeAddress = append(inbound.TunOptions.Inet4RouteExcludeAddress, common.Filter(inbound.TunOptions.RouteExcludeAddress, func(it netip.Prefix) bool {
					return it.Addr().Is4()
				})...)
				inbound.TunOptions.Inet6RouteExcludeAddress = append(inbound.TunOptions.Inet6RouteExcludeAddress, common.Filter(inbound.TunOptions.RouteExcludeAddress, func(it netip.Prefix) bool {
					return it.Addr().Is6()
				})...)
			}
		}
		newInbounds = append(newInbounds, inbound)
	}
	options.Inbounds = newInbounds
}
