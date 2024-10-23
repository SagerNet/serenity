package option

import (
	C "github.com/sagernet/serenity/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-dns"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
)

type _Template struct {
	RawMessage json.RawMessage `json:"-"`
	Name       string          `json:"name,omitempty"`
	Extend     string          `json:"extend,omitempty"`

	// Global

	Log                  *option.LogOptions    `json:"log,omitempty"`
	DomainStrategy       option.DomainStrategy `json:"domain_strategy,omitempty"`
	DomainStrategyLocal  option.DomainStrategy `json:"domain_strategy_local,omitempty"`
	DisableTrafficBypass bool                  `json:"disable_traffic_bypass,omitempty"`
	DisableRuleSet       bool                  `json:"disable_rule_set,omitempty"`
	RemoteResolve        bool                  `json:"remote_resolve,omitempty"`

	// DNS
	Servers        []option.DNSServerOptions `json:"servers,omitempty"`
	DNS            string                    `json:"dns,omitempty"`
	DNSLocal       string                    `json:"dns_local,omitempty"`
	EnableFakeIP   bool                      `json:"enable_fakeip,omitempty"`
	DisableDNSLeak bool                      `json:"disable_dns_leak,omitempty"`
	PreDNSRules    []option.DNSRule          `json:"pre_dns_rules,omitempty"`
	CustomDNSRules []option.DNSRule          `json:"custom_dns_rules,omitempty"`

	// Inbound
	Inbounds           []option.Inbound                              `json:"inbounds,omitempty"`
	AutoRedirect       bool                                          `json:"auto_redirect,omitempty"`
	DisableTUN         bool                                          `json:"disable_tun,omitempty"`
	DisableSystemProxy bool                                          `json:"disable_system_proxy,omitempty"`
	CustomTUN          *TypedMessage[option.TunInboundOptions]       `json:"custom_tun,omitempty"`
	CustomMixed        *TypedMessage[option.HTTPMixedInboundOptions] `json:"custom_mixed,omitempty"`

	// Outbound
	ExtraGroups    []ExtraGroup                    `json:"extra_groups,omitempty"`
	DirectTag      string                          `json:"direct_tag,omitempty"`
	BlockTag       string                          `json:"block_tag,omitempty"`
	DefaultTag     string                          `json:"default_tag,omitempty"`
	URLTestTag     string                          `json:"urltest_tag,omitempty"`
	CustomDirect   *option.DirectOutboundOptions   `json:"custom_direct,omitempty"`
	CustomSelector *option.SelectorOutboundOptions `json:"custom_selector,omitempty"`
	CustomURLTest  *option.URLTestOutboundOptions  `json:"custom_urltest,omitempty"`

	// Route
	DisableDefaultRules bool                   `json:"disable_default_rules,omitempty"`
	PreRules            []option.Rule          `json:"pre_rules,omitempty"`
	CustomRules         []option.Rule          `json:"custom_rules,omitempty"`
	EnableJSDelivr      bool                   `json:"enable_jsdelivr,omitempty"`
	CustomGeoIP         *option.GeoIPOptions   `json:"custom_geoip,omitempty"`
	CustomGeosite       *option.GeositeOptions `json:"custom_geosite,omitempty"`
	CustomRuleSet       []RuleSet              `json:"custom_rule_set,omitempty"`
	PostRuleSet         []RuleSet              `json:"post_rule_set,omitempty"`

	//  Experimental
	DisableCacheFile          bool `json:"disable_cache_file,omitempty"`
	DisableExternalController bool `json:"disable_external_controller,omitempty"`
	DisableClashMode          bool `json:"disable_clash_mode,omitempty"`

	ClashModeLeak   string                                `json:"clash_mode_leak,omitempty"`
	ClashModeRule   string                                `json:"clash_mode_rule,omitempty"`
	ClashModeGlobal string                                `json:"clash_mode_global,omitempty"`
	ClashModeDirect string                                `json:"clash_mode_direct,omitempty"`
	CustomClashAPI  *TypedMessage[option.ClashAPIOptions] `json:"custom_clash_api,omitempty"`

	// Debug
	PProfListen string             `json:"pprof_listen,omitempty"`
	MemoryLimit option.MemoryBytes `json:"memory_limit,omitempty"`
}

type Template _Template

func (t *Template) MarshalJSON() ([]byte, error) {
	return json.Marshal((*_Template)(t))
}

func (t *Template) UnmarshalJSON(content []byte) error {
	err := json.UnmarshalDisallowUnknownFields(content, (*_Template)(t))
	if err != nil {
		return err
	}
	t.RawMessage = content
	return nil
}

type _RuleSet struct {
	Type           string               `json:"type,omitempty"`
	DefaultOptions option.RuleSet       `json:"-"`
	GitHubOptions  GitHubRuleSetOptions `json:"-"`
}

type RuleSet _RuleSet

func (r *RuleSet) MarshalJSON() ([]byte, error) {
	if r.Type == C.RuleSetTypeGitHub {
		return option.MarshallObjects((*_RuleSet)(r), r.GitHubOptions)
	} else {
		return json.Marshal(r.DefaultOptions)
	}
}

func (r *RuleSet) UnmarshalJSON(content []byte) error {
	err := json.Unmarshal(content, (*_RuleSet)(r))
	if err != nil {
		return err
	}
	if r.Type == C.RuleSetTypeGitHub {
		return option.UnmarshallExcluded(content, (*_RuleSet)(r), &r.GitHubOptions)
	} else {
		return option.UnmarshallExcluded(content, (*_RuleSet)(r), &r.DefaultOptions)
	}
}

type GitHubRuleSetOptions struct {
	Repository string                  `json:"repository,omitempty"`
	Path       string                  `json:"path,omitempty"`
	Prefix     string                  `json:"prefix,omitempty"`
	RuleSet    option.Listable[string] `json:"rule_set,omitempty"`
}

func (t Template) DisableIPv6() bool {
	return t.DomainStrategy == option.DomainStrategy(dns.DomainStrategyUseIPv4) && t.DomainStrategyLocal == option.DomainStrategy(dns.DomainStrategyUseIPv4)
}

type ExtraGroup struct {
	Tag                string                          `json:"tag,omitempty"`
	Target             ExtraGroupTarget                `json:"target,omitempty"`
	TagPerSubscription string                          `json:"tag_per_subscription,omitempty"`
	Type               string                          `json:"type,omitempty"`
	Filter             option.Listable[string]         `json:"filter,omitempty"`
	Exclude            option.Listable[string]         `json:"exclude,omitempty"`
	CustomSelector     *option.SelectorOutboundOptions `json:"custom_selector,omitempty"`
	CustomURLTest      *option.URLTestOutboundOptions  `json:"custom_urltest,omitempty"`
}

type ExtraGroupTarget uint8

const (
	ExtraGroupTargetDefault ExtraGroupTarget = iota
	ExtraGroupTargetGlobal
	ExtraGroupTargetSubscription
)

func (t ExtraGroupTarget) String() string {
	switch t {
	case ExtraGroupTargetDefault:
		return "default"
	case ExtraGroupTargetGlobal:
		return "global"
	case ExtraGroupTargetSubscription:
		return "subscription"
	default:
		return "unknown"
	}
}

func (t ExtraGroupTarget) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *ExtraGroupTarget) UnmarshalJSON(bytes []byte) error {
	var stringValue string
	err := json.Unmarshal(bytes, &stringValue)
	if err != nil {
		return err
	}
	switch stringValue {
	case "default":
		*t = ExtraGroupTargetDefault
	case "global":
		*t = ExtraGroupTargetGlobal
	case "subscription":
		*t = ExtraGroupTargetSubscription
	default:
		return E.New("unknown extra group target: ", stringValue)
	}
	return nil
}
