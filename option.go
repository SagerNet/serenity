package serenity

import (
	"bytes"
	"strings"
	"time"

	"github.com/sagernet/serenity/libsubscription"
	"github.com/sagernet/sing-box/common/json"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

const (
	PlatformAndroid = "android"
	PlatformiOS     = "ios"
)

type _Options struct {
	Log            *option.LogOptions        `json:"log,omitempty"`
	Listen         string                    `json:"listen,omitempty"`
	TLS            *option.InboundTLSOptions `json:"tls,omitempty"`
	Subscriptions  []*SubscriptionOptions    `json:"subscriptions,omitempty"`
	Outbounds      []option.Outbound         `json:"outbounds,omitempty"`
	Profiles       []ProfileOptions          `json:"profiles,omitempty"`
	DefaultProfile string                    `json:"default_profile,omitempty"`
}

type UserOptions struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type SubscriptionOptions struct {
	Name             string          `json:"name,omitempty"`
	URL              string          `json:"url,omitempty"`
	UserAgent        string          `json:"user_agent,omitempty"`
	UpdateInterval   option.Duration `json:"update_interval,omitempty"`
	GenerateSelector bool            `json:"generate_selector,omitempty"`
	GenerateURLTest  bool            `json:"generate_url_test,omitempty"`

	LastUpdate  time.Time                `json:"-"`
	ServerCache []libsubscription.Server `json:"-"`
}

type ProfileOptions struct {
	Name               string                  `json:"name,omitempty"`
	Template           string                  `json:"template,omitempty"`
	Config             option.Listable[string] `json:"config,omitempty"`
	ConfigDirectory    option.Listable[string] `json:"config_directory,omitempty"`
	GroupTag           option.Listable[string] `json:"group_tag,omitempty"`
	FilterSubscription option.Listable[string] `json:"filter_subscription,omitempty"`
	FilterOutbound     option.Listable[string] `json:"filter_outbound,omitempty"`
	Authorization      *UserOptions            `json:"authorization,omitempty"`
	Debug              bool                    `json:"debug,omitempty"`
}

type Options _Options

func (o *Options) UnmarshalJSON(content []byte) error {
	decoder := json.NewDecoder(json.NewCommentFilter(bytes.NewReader(content)))
	decoder.DisallowUnknownFields()
	err := decoder.Decode((*_Options)(o))
	if err == nil {
		return nil
	}
	if syntaxError, isSyntaxError := err.(*json.SyntaxError); isSyntaxError {
		prefix := string(content[:syntaxError.Offset])
		row := strings.Count(prefix, "\n") + 1
		column := len(prefix) - strings.LastIndex(prefix, "\n") - 1
		return E.Extend(syntaxError, "row ", row, ", column ", column)
	}
	return err
}
