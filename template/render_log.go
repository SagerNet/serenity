package template

import (
	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/sing-box/option"
)

func (t *Template) renderLog(_ M.Metadata, options *option.Options) error {
	if t.CustomLog != nil {
		options.Log = t.CustomLog
	} else {
		options.Log = &option.LogOptions{
			Level: "info",
		}
	}

	return nil
}
