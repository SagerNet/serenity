package serenity

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sagernet/sing-box/common/badjsonmerge"
	"github.com/sagernet/sing-box/common/badversion"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

type Profile struct {
	name               string
	template           string
	options            option.Options
	groupTag           []string
	filterSubscription []string
	filterOutbound     []string
	authorization      *UserOptions
	debug              bool
}

func NewProfile(options ProfileOptions) (*Profile, error) {
	if options.Template != "" {
		switch options.Template {
		case "default":
		default:
			return nil, E.New("unknown template: ", options.Template)
		}
	} else {
		if len(options.Config) == 0 && len(options.ConfigDirectory) == 0 {
			return nil, E.New("missing configuration")
		}
	}
	mergedConfig, err := readConfigAndMerge(options.Config, options.ConfigDirectory)
	if err != nil {
		return nil, err
	}
	return &Profile{
		name:               options.Name,
		template:           options.Template,
		options:            mergedConfig,
		groupTag:           options.GroupTag,
		filterSubscription: options.FilterSubscription,
		filterOutbound:     options.FilterOutbound,
		authorization:      options.Authorization,
		debug:              options.Debug,
	}, nil
}

func (p *Profile) Name() string {
	return p.name
}

func (p *Profile) GenerateConfig(platform string, version *badversion.Version, outbounds []option.Outbound, subscriptions []*SubscriptionOptions) option.Options {
	options := p.options
	groupTag := p.groupTag
	var template *Profile
	switch p.template {
	case "default":
		template = DefaultTemplate(platform, version, p.debug)
	}
	if template != nil {
		options, _ = badjsonmerge.MergeOptions(options, template.options)
		groupTag = append(groupTag, template.groupTag...)
	}

	if len(p.filterOutbound) > 0 {
		outbounds = common.Filter(outbounds, func(it option.Outbound) bool {
			return common.Contains(p.filterOutbound, it.Tag)
		})
	}
	if len(p.filterSubscription) > 0 {
		subscriptions = common.Filter(subscriptions, func(it *SubscriptionOptions) bool {
			return common.Contains(p.filterSubscription, it.Name)
		})
	}
	if version != nil && badversion.Parse("1.3-beta9").After(*version) {
		outbounds = common.Filter(outbounds, filterH2Mux)
		outbounds = common.Filter(outbounds, filterMuxPadding)
	}
	groupOutbounds := common.Map(common.Filter(options.Outbounds, func(it option.Outbound) bool {
		return common.Contains(groupTag, it.Tag)
	}), func(it option.Outbound) *option.Outbound {
		return &it
	})
	for _, outbound := range outbounds {
		options.Outbounds = append(options.Outbounds, outbound)
		if outbound.Tag != "" {
			for _, groupOutbound := range groupOutbounds {
				switch groupOutbound.Type {
				case C.TypeSelector:
					groupOutbound.SelectorOptions.Outbounds = append(groupOutbound.SelectorOptions.Outbounds, outbound.Tag)
				case C.TypeURLTest:
					groupOutbound.URLTestOptions.Outbounds = append(groupOutbound.URLTestOptions.Outbounds, outbound.Tag)
				}
			}
		}
	}
	for _, subscription := range subscriptions {
		var subscriptionSelector option.Outbound
		subscriptionSelector.Tag = subscription.Name
		subscriptionSelector.Type = C.TypeSelector

		var subscriptionURLTest option.Outbound
		subscriptionURLTest.Tag = subscription.Name + "-auto"
		subscriptionURLTest.Type = C.TypeURLTest

		for _, server := range subscription.ServerCache {
			options.Outbounds = append(options.Outbounds, server.Outbounds...)
			subscriptionSelector.SelectorOptions.Outbounds = append(subscriptionSelector.SelectorOptions.Outbounds, server.Name)
			subscriptionURLTest.URLTestOptions.Outbounds = append(subscriptionURLTest.URLTestOptions.Outbounds, server.Name)
			if !subscription.GenerateSelector {
				for _, groupOutbound := range groupOutbounds {
					switch groupOutbound.Type {
					case C.TypeSelector:
						groupOutbound.SelectorOptions.Outbounds = append(groupOutbound.SelectorOptions.Outbounds, server.Name)
					case C.TypeURLTest:
						groupOutbound.URLTestOptions.Outbounds = append(groupOutbound.URLTestOptions.Outbounds, server.Name)
					}
				}
			}
		}
		if subscription.GenerateSelector {
			options.Outbounds = append(options.Outbounds, subscriptionSelector)
			for _, groupOutbound := range groupOutbounds {
				switch groupOutbound.Type {
				case C.TypeSelector:
					groupOutbound.SelectorOptions.Outbounds = append(groupOutbound.SelectorOptions.Outbounds, subscriptionSelector.Tag)
					if groupOutbound.SelectorOptions.Default == "" {
						groupOutbound.SelectorOptions.Default = subscriptionSelector.Tag
					}
				case C.TypeURLTest:
					groupOutbound.URLTestOptions.Outbounds = append(groupOutbound.URLTestOptions.Outbounds, subscriptionSelector.Tag)
				}
			}
		}
		if subscription.GenerateURLTest {
			options.Outbounds = append(options.Outbounds, subscriptionURLTest)
			for _, groupOutbound := range groupOutbounds {
				switch groupOutbound.Type {
				case C.TypeSelector:
					groupOutbound.SelectorOptions.Outbounds = append(groupOutbound.SelectorOptions.Outbounds, subscriptionURLTest.Tag)
					if groupOutbound.SelectorOptions.Default == "" {
						groupOutbound.SelectorOptions.Default = subscriptionURLTest.Tag
					}
				}
			}
		}
	}
	for _, groupOutbound := range groupOutbounds {
		for i, outbound := range options.Outbounds {
			if outbound.Tag == groupOutbound.Tag {
				options.Outbounds[i] = *groupOutbound
				break
			}
		}
	}
	return options
}

func (p *Profile) CheckBasicAuthorization(request *http.Request) bool {
	options := p.authorization
	if options == nil || options.Username == "" {
		return true
	}
	header := request.Header.Get("Authorization")
	basic, encoded, found := strings.Cut(header, " ")
	if !found || basic != "Basic" {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return false
	}
	username, password, found := strings.Cut(string(decoded), ":")
	if !found {
		return false
	}
	return username == options.Username && password == options.Password
}

type profileOptionsEntry struct {
	content []byte
	path    string
	options option.Options
}

func readConfigAt(path string) (*profileOptionsEntry, error) {
	var (
		configContent []byte
		err           error
	)
	if path == "stdin" {
		configContent, err = io.ReadAll(os.Stdin)
	} else {
		configContent, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, E.Cause(err, "read config at ", path)
	}
	var options option.Options
	err = options.UnmarshalJSON(configContent)
	if err != nil {
		return nil, E.Cause(err, "decode config at ", path)
	}
	return &profileOptionsEntry{
		content: configContent,
		path:    path,
		options: options,
	}, nil
}

func readConfig(configPath []string, configDirectory []string) ([]*profileOptionsEntry, error) {
	var optionsList []*profileOptionsEntry
	for _, path := range configPath {
		optionsEntry, err := readConfigAt(path)
		if err != nil {
			return nil, err
		}
		optionsList = append(optionsList, optionsEntry)
	}
	for _, directory := range configDirectory {
		entries, err := os.ReadDir(directory)
		if err != nil {
			return nil, E.Cause(err, "read config directory at ", directory)
		}
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".json") || entry.IsDir() {
				continue
			}
			optionsEntry, err := readConfigAt(filepath.Join(directory, entry.Name()))
			if err != nil {
				return nil, err
			}
			optionsList = append(optionsList, optionsEntry)
		}
	}
	sort.Slice(optionsList, func(i, j int) bool {
		return optionsList[i].path < optionsList[j].path
	})
	return optionsList, nil
}

func readConfigAndMerge(configPath []string, configDirectory []string) (option.Options, error) {
	optionsList, err := readConfig(configPath, configDirectory)
	if err != nil {
		return option.Options{}, err
	}
	if len(optionsList) == 1 {
		return optionsList[0].options, nil
	}
	var mergedOptions option.Options
	for _, options := range optionsList {
		mergedOptions, err = badjsonmerge.MergeOptions(options.options, mergedOptions)
		if err != nil {
			return option.Options{}, E.Cause(err, "merge config at ", options.path)
		}
	}
	return mergedOptions, nil
}
