package template

import (
	"context"
	"regexp"

	"github.com/sagernet/serenity/option"
	C "github.com/sagernet/sing-box/constant"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
	"github.com/sagernet/sing/common/logger"
)

type Manager struct {
	ctx       context.Context
	logger    logger.Logger
	templates []*Template
}

func extendTemplate(rawTemplates []option.Template, root, current option.Template) (option.Template, error) {
	if current.Extend == "" {
		return current, nil
	} else if root.Name == current.Extend {
		return option.Template{}, E.New("initialize template[", current.Name, "]: circular extend detected: ", current.Extend)
	}
	var next option.Template
	for _, it := range rawTemplates {
		if it.Name == current.Extend {
			next = it
			break
		}
	}
	if next.Name == "" {
		return option.Template{}, E.New("initialize template[", current.Name, "]: extended template not found: ", current.Extend)
	}
	if next.Extend != "" {
		newNext, err := extendTemplate(rawTemplates, root, next)
		if err != nil {
			return option.Template{}, E.Cause(err, next.Extend)
		}
		next = newNext
	}
	newRawTemplate, err := badjson.MergeJSON(next.RawMessage, current.RawMessage)
	if err != nil {
		return option.Template{}, E.Cause(err, "initialize template[", current.Name, "]: merge extended template: ", current.Extend)
	}
	newTemplate, err := json.UnmarshalExtended[option.Template](newRawTemplate)
	if err != nil {
		return option.Template{}, E.Cause(err, "initialize template[", current.Name, "]: unmarshal extended template: ", current.Extend)
	}
	newTemplate.RawMessage = newRawTemplate
	return newTemplate, nil
}

func NewManager(ctx context.Context, logger logger.Logger, rawTemplates []option.Template) (*Manager, error) {
	var templates []*Template
	for templateIndex, template := range rawTemplates {
		if template.Name == "" {
			return nil, E.New("initialize template[", templateIndex, "]: missing name")
		}
		if template.Extend != "" {
			newTemplate, err := extendTemplate(rawTemplates, template, template)
			if err != nil {
				return nil, err
			}
			template = newTemplate
		}
		var groups []*ExtraGroup
		for groupIndex, group := range template.ExtraGroups {
			if group.Tag == "" {
				return nil, E.New("initialize template[", template.Name, "]: extra_group[", groupIndex, "]: missing tag")
			}
			switch group.Type {
			case C.TypeSelector, C.TypeURLTest:
			case "":
				return nil, E.New("initialize template[", template.Name, "]: extra_group[", group.Tag, "]: missing type")
			default:
				return nil, E.New("initialize template[", template.Name, "]: extra_group[", group.Tag, "]: invalid group type: ", group.Type)
			}
			var (
				filter  []*regexp.Regexp
				exclude []*regexp.Regexp
			)
			for filterIndex, it := range group.Filter {
				regex, err := regexp.Compile(it)
				if err != nil {
					return nil, E.Cause(err, "initialize template[", template.Name, "]: parse extra_group[", group.Tag, "]: parse filter[", filterIndex, "]: ", it)
				}
				filter = append(filter, regex)
			}
			for excludeIndex, it := range group.Exclude {
				regex, err := regexp.Compile(it)
				if err != nil {
					return nil, E.Cause(err, "initialize template[", template.Name, "]: parse extra_group[", group.Tag, "]: parse exclude[", excludeIndex, "]: ", it)
				}
				exclude = append(exclude, regex)
			}
			groups = append(groups, &ExtraGroup{
				ExtraGroup: group,
				filter:     filter,
				exclude:    exclude,
			})
		}
		templates = append(templates, &Template{
			Template: template,
			groups:   groups,
		})
	}
	return &Manager{
		ctx:       ctx,
		logger:    logger,
		templates: templates,
	}, nil
}

func (m *Manager) TemplateByName(name string) *Template {
	for _, template := range m.templates {
		if template.Name == name {
			return template
		}
	}
	return nil
}
