package template

import (
	"regexp"

	M "github.com/sagernet/serenity/common/metadata"
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
	DefaultDefaultTag = "Default"
	DefaultDirectTag  = "direct"
	BlockTag          = "block"
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

func (t *Template) Render(metadata M.Metadata, profileName string, outbounds [][]boxOption.Outbound, subscriptions []*subscription.Subscription) (*boxOption.Options, error) {
	var options boxOption.Options
	err := t.renderDNS(metadata, &options)
	if err != nil {
		return nil, E.Cause(err, "render dns")
	}
	err = t.renderInbounds(metadata, &options)
	if err != nil {
		return nil, E.Cause(err, "render inbounds")
	}
	err = t.renderOutbounds(metadata, &options, outbounds, subscriptions)
	if err != nil {
		return nil, E.Cause(err, "render outbounds")
	}
	err = t.renderRoute(metadata, &options)
	if err != nil {
		return nil, E.Cause(err, "render route")
	}
	err = t.renderExperimental(metadata, &options, profileName)
	if err != nil {
		return nil, E.Cause(err, "render experimental")
	}
	filter.Filter(metadata, &options)
	return &options, nil
}
