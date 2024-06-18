package filter

import (
	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func init() {
	filters = append(filters, filter170)
}

func filter170(metadata metadata.Metadata, options *option.Options) {
	if metadata.Version == nil || metadata.Version.GreaterThanOrEqual(semver.Version{Major: 1, Minor: 7}) {
		return
	}
	newInbounds := make([]option.Inbound, 0, len(options.Inbounds))
	for _, inbound := range options.Inbounds {
		switch inbound.Type {
		case C.TypeTun:
			inbound.TunOptions = filter170Tun(inbound.TunOptions)
			inbound.TunOptions.InboundOptions = filter170InboundOptions(inbound.TunOptions.InboundOptions)
		case C.TypeRedirect:
			inbound.RedirectOptions.InboundOptions = filter170InboundOptions(inbound.RedirectOptions.InboundOptions)
		case C.TypeTProxy:
			inbound.TProxyOptions.InboundOptions = filter170InboundOptions(inbound.TProxyOptions.InboundOptions)
		case C.TypeDirect:
			inbound.DirectOptions.InboundOptions = filter170InboundOptions(inbound.DirectOptions.InboundOptions)
		case C.TypeSOCKS:
			inbound.SocksOptions.InboundOptions = filter170InboundOptions(inbound.SocksOptions.InboundOptions)
		case C.TypeHTTP:
			inbound.HTTPOptions.InboundOptions = filter170InboundOptions(inbound.HTTPOptions.InboundOptions)
		case C.TypeMixed:
			inbound.MixedOptions.InboundOptions = filter170InboundOptions(inbound.MixedOptions.InboundOptions)
		case C.TypeShadowsocks:
			inbound.ShadowsocksOptions.InboundOptions = filter170InboundOptions(inbound.ShadowsocksOptions.InboundOptions)
			inbound.ShadowsocksOptions.Multiplex = nil
		case C.TypeVMess:
			inbound.VMessOptions.InboundOptions = filter170InboundOptions(inbound.VMessOptions.InboundOptions)
			inbound.VMessOptions.Multiplex = nil
		case C.TypeTrojan:
			inbound.TrojanOptions.InboundOptions = filter170InboundOptions(inbound.TrojanOptions.InboundOptions)
			inbound.TrojanOptions.Multiplex = nil
		case C.TypeNaive:
			inbound.NaiveOptions.InboundOptions = filter170InboundOptions(inbound.NaiveOptions.InboundOptions)
		case C.TypeHysteria:
			inbound.HysteriaOptions.InboundOptions = filter170InboundOptions(inbound.HysteriaOptions.InboundOptions)
		case C.TypeShadowTLS:
			inbound.ShadowTLSOptions.InboundOptions = filter170InboundOptions(inbound.ShadowTLSOptions.InboundOptions)
		case C.TypeVLESS:
			inbound.VLESSOptions.InboundOptions = filter170InboundOptions(inbound.VLESSOptions.InboundOptions)
			inbound.VLESSOptions.Multiplex = nil
		case C.TypeTUIC:
			inbound.TUICOptions.InboundOptions = filter170InboundOptions(inbound.TUICOptions.InboundOptions)
		case C.TypeHysteria2:
			inbound.Hysteria2Options.InboundOptions = filter170InboundOptions(inbound.Hysteria2Options.InboundOptions)
		default:
			continue
		}
		newInbounds = append(newInbounds, inbound)
	}
	options.Inbounds = newInbounds
	if options.Route != nil {
		options.Route.Rules = common.Filter(options.Route.Rules, filter170Rule)
	}
	if options.DNS != nil {
		options.DNS.Rules = common.Filter(options.DNS.Rules, filter170DNSRule)
	}
}

//nolint:staticcheck
//goland:noinspection GoDeprecation
func filter170Tun(options option.TunInboundOptions) option.TunInboundOptions {
	options.Inet4RouteExcludeAddress = nil
	options.Inet6RouteExcludeAddress = nil
	return options
}

func filter170InboundOptions(options option.InboundOptions) option.InboundOptions {
	options.UDPDisableDomainUnmapping = false
	return options
}

func filter170Rule(it option.Rule) bool {
	return !hasRule([]option.Rule{it}, isWIFIRule)
}

func filter170DNSRule(it option.DNSRule) bool {
	return !hasDNSRule([]option.DNSRule{it}, isWIFIDNSRule)
}
