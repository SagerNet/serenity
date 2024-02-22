package server

import (
	"context"
	"regexp"

	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/option"
	"github.com/sagernet/serenity/subscription"
	"github.com/sagernet/serenity/template"
	boxOption "github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/logger"
)

type ProfileManager struct {
	ctx            context.Context
	logger         logger.Logger
	subscription   *subscription.Manager
	outbounds      [][]boxOption.Outbound
	profiles       []*Profile
	defaultProfile *Profile
}

type Profile struct {
	option.Profile
	manager              *ProfileManager
	template             *template.Template
	templateForPlatform  map[metadata.Platform]*template.Template
	templateForUserAgent map[*regexp.Regexp]*template.Template
	groups               []ExtraGroup
}

type ExtraGroup struct {
	option.ExtraGroup
	filterRegex []*regexp.Regexp
}

func NewProfileManager(
	ctx context.Context,
	logger logger.Logger,
	subscriptionManager *subscription.Manager,
	templateManager *template.Manager,
	outbounds [][]boxOption.Outbound,
	rawProfiles []option.Profile,
) (*ProfileManager, error) {
	manager := &ProfileManager{
		ctx:          ctx,
		logger:       logger,
		subscription: subscriptionManager,
		outbounds:    outbounds,
	}
	for profileIndex, profile := range rawProfiles {
		if profile.Name == "" {
			return nil, E.New("initialize profile[", profileIndex, "]: missing name")
		}
		var (
			defaultTemplate      *template.Template
			templateForPlatform  = make(map[metadata.Platform]*template.Template)
			templateForUserAgent = make(map[*regexp.Regexp]*template.Template)
		)
		if profile.Template != "" {
			defaultTemplate = templateManager.TemplateByName(profile.Template)
			if defaultTemplate == nil {
				return nil, E.New("initialize profile[", profile.Name, "]: template not found: ", profile.Template)
			}
		} else {
			defaultTemplate = template.Default
		}
		if profile.TemplateForPlatform != nil {
			for templateIndex, entry := range profile.TemplateForPlatform.Entries() {
				platform, err := metadata.ParsePlatform(entry.Key)
				if err != nil {
					return nil, E.Cause(err, "initialize profile[", profile.Name, "]: parse template_for_platform[", templateIndex, "]")
				}
				customTemplate := templateManager.TemplateByName(entry.Value)
				if customTemplate == nil {
					return nil, E.New("initialize profile[", profile.Name, "]: parse template_for_platform[", entry.Key, "]: template not found: ", entry.Value)
				}
				templateForPlatform[platform] = customTemplate
			}
		}
		if profile.TemplateForUserAgent != nil {
			for templateIndex, entry := range profile.TemplateForUserAgent.Entries() {
				regex, err := regexp.Compile(entry.Key)
				if err != nil {
					return nil, E.Cause(err, "initialize profile[", profile.Name, "]: parse template_for_user_agent[", templateIndex, "]")
				}
				customTemplate := templateManager.TemplateByName(entry.Value)
				if customTemplate == nil {
					return nil, E.New("initialize profile[", profile.Name, "]: parse template_for_user_agent[", entry.Key, "]: template not found: ", entry.Value)
				}
				templateForUserAgent[regex] = customTemplate
			}
		}
		manager.profiles = append(manager.profiles, &Profile{
			Profile:              profile,
			manager:              manager,
			template:             defaultTemplate,
			templateForPlatform:  templateForPlatform,
			templateForUserAgent: templateForUserAgent,
		})
	}
	if len(manager.profiles) > 0 {
		manager.defaultProfile = manager.profiles[0]
	}
	return manager, nil
}

func (m *ProfileManager) ProfileByName(name string) *Profile {
	for _, it := range m.profiles {
		if it.Name == name {
			return it
		}
	}
	return nil
}

func (m *ProfileManager) DefaultProfile() *Profile {
	return m.defaultProfile
}

func (p *Profile) Render(metadata metadata.Metadata) (*boxOption.Options, error) {
	selectedTemplate, loaded := p.templateForPlatform[metadata.Platform]
	if !loaded {
		for regex, it := range p.templateForUserAgent {
			if regex.MatchString(metadata.UserAgent) {
				selectedTemplate = it
				break
			}
		}
	}
	if selectedTemplate == nil {
		selectedTemplate = p.template
	}
	outbounds := common.Filter(p.manager.outbounds, func(it []boxOption.Outbound) bool {
		return common.Contains(p.Outbound, it[0].Tag)
	})
	subscriptions := common.Filter(p.manager.subscription.Subscriptions(), func(it *subscription.Subscription) bool {
		return common.Contains(p.Subscription, it.Name)
	})
	options, err := selectedTemplate.Render(metadata, p.Name, outbounds, subscriptions)
	if err != nil {
		return nil, err
	}
	options, err = badjson.Omitempty(options)
	if err != nil {
		return nil, E.Cause(err, "omitempty")
	}
	return options, nil
}
