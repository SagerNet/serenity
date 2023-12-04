package filter

import (
	"github.com/sagernet/serenity/common/metadata"
	boxOption "github.com/sagernet/sing-box/option"
)

type OptionsFilter func(metadata metadata.Metadata, options *boxOption.Options)

var filters []OptionsFilter

func Filter(metadata metadata.Metadata, options *boxOption.Options) {
	for _, filter := range filters {
		filter(metadata, options)
	}
}
