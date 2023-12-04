package template

import (
	"context"
	"regexp"

	"github.com/sagernet/serenity/option"
	C "github.com/sagernet/sing-box/constant"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"
)

type Manager struct {
	ctx       context.Context
	logger    logger.Logger
	templates []*Template
}

func NewManager(ctx context.Context, logger logger.Logger, rawTemplates []option.Template) (*Manager, error) {
	var templates []*Template
	for templateIndex, template := range rawTemplates {
		if template.Name == "" {
			return nil, E.New("initialize template[", templateIndex, "]: missing name")
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
