package option

import (
	"bytes"
	"time"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
)

type _Options struct {
	RawMessage json.RawMessage           `json:"-"`
	Log        *option.LogOptions        `json:"log,omitempty"`
	Listen     string                    `json:"listen,omitempty"`
	TLS        *option.InboundTLSOptions `json:"tls,omitempty"`
	CacheFile  string                    `json:"cache_file,omitempty"`

	Outbounds     []option.Listable[option.Outbound] `json:"outbounds,omitempty"`
	Subscriptions []Subscription                     `json:"subscriptions,omitempty"`
	Templates     []Template                         `json:"templates,omitempty"`
	Profiles      []Profile                          `json:"profiles,omitempty"`
	Users         []User                             `json:"users,omitempty"`
}

type Options _Options

func (o *Options) UnmarshalJSON(content []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err := decoder.Decode((*_Options)(o))
	if err != nil {
		return err
	}
	o.RawMessage = content
	return nil
}

type User struct {
	Name           string                  `json:"name,omitempty"`
	Password       string                  `json:"password,omitempty"`
	Profile        option.Listable[string] `json:"profile,omitempty"`
	DefaultProfile string                  `json:"default_profile,omitempty"`
}

const (
	DefaultSubscriptionUpdateInterval = 1 * time.Hour
)

type Subscription struct {
	Name             string                                  `json:"name,omitempty"`
	URL              string                                  `json:"url,omitempty"`
	UserAgent        string                                  `json:"user_agent,omitempty"`
	UpdateInterval   option.Duration                         `json:"update_interval,omitempty"`
	Process          option.Listable[OutboundProcessOptions] `json:"process,omitempty"`
	DeDuplication    bool                                    `json:"deduplication,omitempty"`
	GenerateSelector bool                                    `json:"generate_selector,omitempty"`
	GenerateURLTest  bool                                    `json:"generate_urltest,omitempty"`
	URLTestTagSuffix string                                  `json:"urltest_suffix,omitempty"`
	CustomSelector   *option.SelectorOutboundOptions         `json:"custom_selector,omitempty"`
	CustomURLTest    *option.URLTestOutboundOptions          `json:"custom_urltest,omitempty"`
}

type OutboundProcessOptions struct {
	Filter           option.Listable[string]           `json:"filter,omitempty"`
	Exclude          option.Listable[string]           `json:"exclude,omitempty"`
	FilterType       option.Listable[string]           `json:"filter_type,omitempty"`
	ExcludeType      option.Listable[string]           `json:"exclude_type,omitempty"`
	Invert           bool                              `json:"invert,omitempty"`
	Remove           bool                              `json:"remove,omitempty"`
	Rename           *badjson.TypedMap[string, string] `json:"rename,omitempty"`
	RemoveEmoji      bool                              `json:"remove_emoji,omitempty"`
	RewriteMultiplex *option.OutboundMultiplexOptions  `json:"rewrite_multiplex,omitempty"`
}

type Profile struct {
	Name                 string                            `json:"name,omitempty"`
	Template             string                            `json:"template,omitempty"`
	TemplateForPlatform  *badjson.TypedMap[string, string] `json:"template_for_platform,omitempty"`
	TemplateForUserAgent *badjson.TypedMap[string, string] `json:"template_for_user_agent,omitempty"`
	Outbound             option.Listable[string]           `json:"outbound,omitempty"`
	Subscription         option.Listable[string]           `json:"subscription,omitempty"`
}
