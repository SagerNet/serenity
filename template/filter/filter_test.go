package filter

import (
	"testing"

	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"

	"github.com/stretchr/testify/require"
)

func TestFilter1100(t *testing.T) {
	options := &option.Options{
		DNS: &option.DNSOptions{
			Rules: []option.DNSRule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultDNSRule{
						RuleSet: []string{"test"},
						Server:  "test",
					},
				},
			},
		},
		Route: &option.RouteOptions{
			Rules: []option.Rule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						RuleSet:  []string{"test"},
						Outbound: "test",
					},
				},
			},
			RuleSet: []option.RuleSet{
				{
					Type: C.RuleSetTypeInline,
					Tag:  "test",
					InlineOptions: option.PlainRuleSet{
						Rules: []option.HeadlessRule{
							{
								Type: C.RuleTypeDefault,
								DefaultOptions: option.DefaultHeadlessRule{
									Domain: []string{"example.com"},
								},
							},
						},
					},
				},
			},
		},
	}
	err := filter1100(metadata.Metadata{Version: &semver.Version{Major: 1, Minor: 9, Patch: 3}}, options)
	require.NoError(t, err)
	require.Equal(t, options, &option.Options{
		DNS: &option.DNSOptions{
			Rules: []option.DNSRule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultDNSRule{
						Domain: []string{"example.com"},
						Server: "test",
					},
				},
			},
		},
		Route: &option.RouteOptions{
			Rules: []option.Rule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						Domain:   []string{"example.com"},
						Outbound: "test",
					},
				},
			},
		},
	})
}
