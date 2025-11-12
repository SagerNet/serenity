package template

import (
	"context"
	"regexp"

	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	"github.com/sagernet/serenity/option"
	"github.com/sagernet/serenity/subscription"
	"github.com/sagernet/serenity/template/filter"
	boxOption "github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

const (
	DefaultMixedPort  = 8080
	DNSDefaultTag     = "default"
	DNSLocalTag       = "local"
	DNSLocalSetupTag  = "local_setup"
	DNSFakeIPTag      = "remote"
	DefaultDNS        = "tls://8.8.8.8"
	DefaultDNSLocal   = "https://223.5.5.5/dns-query"
	DefaultDefaultTag = "default"
	DefaultDirectTag  = "direct"
	DefaultBlockTag   = "block"
	DNSTag            = "dns"
	DefaultURLTestTag = "URLTest"
)

var Default = new(Template)

type Template struct {
	option.Template
	groups []*ExtraGroup
}

type ExtraGroup struct {
	option.ExtraGroup
	filter  []*regexp.Regexp
	exclude []*regexp.Regexp
}

func (t *Template) Render(ctx context.Context, metadata M.Metadata, profileName string, endpoints []boxOption.Endpoint, outbounds [][]boxOption.Outbound, subscriptions []*subscription.Subscription) (*boxOption.Options, error) {
	var options boxOption.Options
	options.Log = t.Log
	if metadata.Version != nil && metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.12.0-alpha.1")) {
		options.Endpoints = endpoints
	}
	err := t.renderDNS(ctx, metadata, &options)
	if err != nil {
		return nil, E.Cause(err, "render dns")
	}
	err = t.renderRoute(metadata, &options)
	if err != nil {
		return nil, E.Cause(err, "render route")
	}
	err = t.renderInbounds(metadata, &options)
	if err != nil {
		return nil, E.Cause(err, "render inbounds")
	}
	err = t.renderOutbounds(metadata, &options, outbounds, subscriptions)
	if err != nil {
		return nil, E.Cause(err, "render outbounds")
	}
	err = t.renderExperimental(ctx, metadata, &options, profileName)
	if err != nil {
		return nil, E.Cause(err, "render experimental")
	}
	err = filter.Filter(metadata, &options)
	if err != nil {
		return nil, E.Cause(err, "filter options")
	}
	return &options, nil
}
