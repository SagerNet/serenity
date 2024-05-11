package option

import (
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/serenity/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-dns"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
)

type Template struct {
	Name   string `json:"name,omitempty"`
	Extend string `json:"extend,omitempty"`

	// Global

	Log                  *option.LogOptions    `json:"log,omitempty"`
	DomainStrategy       option.DomainStrategy `json:"domain_strategy,omitempty"`
	DisableTrafficBypass bool                  `json:"disable_traffic_bypass,omitempty"`
	DisableRuleSet       bool                  `json:"disable_rule_set,omitempty"`
	RemoteResolve        bool                  `json:"remote_resolve,omitempty"`

	// DNS
	DNSDefault     string           `json:"dns_default,omitempty"`
	DNSLocal       string           `json:"dns_local,omitempty"`
	EnableFakeIP   bool             `json:"enable_fakeip,omitempty"`
	DisableDNSLeak bool             `json:"disable_dns_leak,omitempty"`
	PreDNSRules    []option.DNSRule `json:"pre_dns_rules,omitempty"`
	CustomDNSRules []option.DNSRule `json:"custom_dns_rules,omitempty"`

	// Inbound
	Inbounds           []option.Inbound                              `json:"inbounds,omitempty"`
	DisableTUN         bool                                          `json:"disable_tun,omitempty"`
	DisableSystemProxy bool                                          `json:"disable_system_proxy,omitempty"`
	CustomTUN          *TypedMessage[option.TunInboundOptions]       `json:"custom_tun,omitempty"`
	CustomMixed        *TypedMessage[option.HTTPMixedInboundOptions] `json:"custom_mixed,omitempty"`

	// Outbound
	ExtraGroups           []ExtraGroup                    `json:"extra_groups,omitempty"`
	GenerateGlobalURLTest bool                            `json:"generate_global_urltest,omitempty"`
	DirectTag             string                          `json:"direct_tag,omitempty"`
	BlockTag              string                          `json:"block_tag,omitempty"`
	DefaultTag            string                          `json:"default_tag,omitempty"`
	URLTestTag            string                          `json:"urltest_tag,omitempty"`
	CustomDirect          *option.DirectOutboundOptions   `json:"custom_direct,omitempty"`
	CustomSelector        *option.SelectorOutboundOptions `json:"custom_selector,omitempty"`
	CustomURLTest         *option.URLTestOutboundOptions  `json:"custom_urltest,omitempty"`

	// Route
	DisableDefaultRules           bool                                            `json:"disable_default_rules,omitempty"`
	PreRules                      []option.Rule                                   `json:"pre_rules,omitempty"`
	CustomRules                   []option.Rule                                   `json:"custom_rules,omitempty"`
	CustomRulesForVersionLessThan badjson.TypedMap[semver.Version, []option.Rule] `json:"custom_rules_for_version_less_than,omitempty"`
	EnableJSDelivr                bool                                            `json:"enable_jsdelivr,omitempty"`
	CustomGeoIP                   *option.GeoIPOptions                            `json:"custom_geoip,omitempty"`
	CustomGeosite                 *option.GeositeOptions                          `json:"custom_geosite,omitempty"`
	CustomRuleSet                 []RuleSet                                       `json:"custom_rule_set,omitempty"`
	PostRuleSet                   []RuleSet                                       `json:"post_rule_set,omitempty"`

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

type _RuleSet struct {
	Type           string               `json:"type,omitempty"`
	DefaultOptions option.RuleSet       `json:"-"`
	GitHubOptions  GitHubRuleSetOptions `json:"-"`
}

type RuleSet _RuleSet

func (r *RuleSet) RawOptions() (any, error) {
	switch r.Type {
	case C.RuleSetTypeDefault, "":
		r.Type = ""
		return &r.DefaultOptions, nil
	case C.RuleSetTypeGitHub:
		return &r.GitHubOptions, nil
	default:
		return nil, E.New("unknown rule set type", r.Type)
	}
}

func (r *RuleSet) MarshalJSON() ([]byte, error) {
	rawOptions, err := r.RawOptions()
	if err != nil {
		return nil, err
	}
	return option.MarshallObjects((*_RuleSet)(r), rawOptions)
}

func (r *RuleSet) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_RuleSet)(r))
	if err != nil {
		return err
	}
	rawOptions, err := r.RawOptions()
	if err != nil {
		return err
	}
	return option.UnmarshallExcluded(bytes, (*_RuleSet)(r), rawOptions)
}

type GitHubRuleSetOptions struct {
	Owner   string                  `json:"owner,omitempty"`
	Repo    string                  `json:"repo,omitempty"`
	Branch  string                  `json:"branch,omitempty"`
	RuleSet option.Listable[string] `json:"rule_set,omitempty"`
}

func (t Template) DisableIPv6() bool {
	return t.DomainStrategy == option.DomainStrategy(dns.DomainStrategyUseIPv4)
}

type ExtraGroup struct {
	Tag            string                          `json:"tag,omitempty"`
	Type           string                          `json:"type,omitempty"`
	Filter         option.Listable[string]         `json:"filter,omitempty"`
	Exclude        option.Listable[string]         `json:"exclude,omitempty"`
	CustomSelector *option.SelectorOutboundOptions `json:"custom_selector,omitempty"`
	CustomURLTest  *option.URLTestOutboundOptions  `json:"custom_urltest,omitempty"`
}
