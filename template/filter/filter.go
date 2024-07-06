package filter

import (
	"github.com/sagernet/serenity/common/metadata"
	boxOption "github.com/sagernet/sing-box/option"
)

type OptionsFilter func(metadata metadata.Metadata, options *boxOption.Options) error

var filters []OptionsFilter

func Filter(metadata metadata.Metadata, options *boxOption.Options) error {
	for _, filter := range filters {
		err := filter(metadata, options)
		if err != nil {
			return err
		}
	}
	return nil
}
