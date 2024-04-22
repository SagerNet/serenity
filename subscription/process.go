package subscription

import (
	"regexp"
	"strings"

	"github.com/sagernet/serenity/option"
	C "github.com/sagernet/sing-box/constant"
	boxOption "github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

type ProcessOptions struct {
	option.OutboundProcessOptions
	filter  []*regexp.Regexp
	exclude []*regexp.Regexp
	rename  []*Rename
}

type Rename struct {
	From *regexp.Regexp
	To   string
}

func NewProcessOptions(options option.OutboundProcessOptions) (*ProcessOptions, error) {
	var (
		filter  []*regexp.Regexp
		exclude []*regexp.Regexp
		rename  []*Rename
	)
	for regexIndex, it := range options.Filter {
		regex, err := regexp.Compile(it)
		if err != nil {
			return nil, E.Cause(err, "parse filter[", regexIndex, "]")
		}
		filter = append(filter, regex)
	}
	for regexIndex, it := range options.Exclude {
		regex, err := regexp.Compile(it)
		if err != nil {
			return nil, E.Cause(err, "parse exclude[", regexIndex, "]")
		}
		exclude = append(exclude, regex)
	}
	if options.Rename != nil {
		for renameIndex, entry := range options.Rename.Entries() {
			regex, err := regexp.Compile(entry.Key)
			if err != nil {
				return nil, E.Cause(err, "parse rename[", renameIndex, "]: parse ", entry.Key)
			}
			rename = append(rename, &Rename{
				From: regex,
				To:   entry.Value,
			})
		}
	}
	return &ProcessOptions{
		OutboundProcessOptions: options,
		filter:                 filter,
		exclude:                exclude,
		rename:                 rename,
	}, nil
}

func (o *ProcessOptions) Process(outbounds []boxOption.Outbound) []boxOption.Outbound {
	newOutbounds := make([]boxOption.Outbound, 0, len(outbounds))
	renameResult := make(map[string]string)
	for _, outbound := range outbounds {
		if len(o.filter) > 0 {
			if !common.Any(o.filter, func(it *regexp.Regexp) bool {
				return it.MatchString(outbound.Tag)
			}) {
				continue
			}
		}
		if len(o.FilterOutboundType) > 0 {
			if !common.Contains(o.FilterOutboundType, outbound.Type) {
				continue
			}
		}
		if len(o.exclude) > 0 {
			if common.Any(o.exclude, func(it *regexp.Regexp) bool {
				return it.MatchString(outbound.Tag)
			}) {
				continue
			}
		}
		if len(o.ExcludeOutboundType) > 0 {
			if common.Contains(o.ExcludeOutboundType, outbound.Type) {
				continue
			}
		}
		originTag := outbound.Tag
		if len(o.rename) > 0 {
			for _, rename := range o.rename {
				outbound.Tag = rename.From.ReplaceAllString(outbound.Tag, rename.To)
			}
		}
		if o.RemoveEmoji {
			outbound.Tag = removeEmojis(outbound.Tag)
		}
		outbound.Tag = strings.TrimSpace(outbound.Tag)
		if originTag != outbound.Tag {
			renameResult[originTag] = outbound.Tag
		}
		if o.RewriteMultiplex != nil {
			switch outbound.Type {
			case C.TypeShadowsocks:
				outbound.ShadowsocksOptions.Multiplex = o.RewriteMultiplex
			case C.TypeTrojan:
				outbound.TrojanOptions.Multiplex = o.RewriteMultiplex
			case C.TypeVMess:
				outbound.VMessOptions.Multiplex = o.RewriteMultiplex
			case C.TypeVLESS:
				outbound.VLESSOptions.Multiplex = o.RewriteMultiplex
			}
		}
		newOutbounds = append(newOutbounds, outbound)
	}
	if len(renameResult) > 0 {
		for i, outbound := range newOutbounds {
			rawOptions, err := outbound.RawOptions()
			if err != nil {
				continue
			}
			if dialerOptionsWrapper, containsDialerOptions := rawOptions.(boxOption.DialerOptionsWrapper); containsDialerOptions {
				dialerOptions := dialerOptionsWrapper.TakeDialerOptions()
				if dialerOptions.Detour == "" {
					continue
				}
				newTag, loaded := renameResult[dialerOptions.Detour]
				if !loaded {
					continue
				}
				dialerOptions.Detour = newTag
				dialerOptionsWrapper.ReplaceDialerOptions(dialerOptions)
				newOutbounds[i] = outbound
			}
		}
	}
	return newOutbounds
}

func removeEmojis(s string) string {
	var runes []rune
	for _, r := range s {
		if !(r >= 0x1F600 && r <= 0x1F64F || // Emoticons
			r >= 0x1F300 && r <= 0x1F5FF || // Symbols & Pictographs
			r >= 0x1F680 && r <= 0x1F6FF || // Transport & Map Symbols
			r >= 0x1F1E0 && r <= 0x1F1FF || // Flags
			r >= 0x2600 && r <= 0x26FF || // Misc symbols
			r >= 0x2700 && r <= 0x27BF || // Dingbats
			r >= 0xFE00 && r <= 0xFE0F || // Variation Selectors
			r >= 0x1F900 && r <= 0x1F9FF || // Supplemental Symbols and Pictographs
			r >= 0x1F018 && r <= 0x1F270 || // Various asian characters
			r >= 0x238C && r <= 0x2454 || // Misc items
			r >= 0x20D0 && r <= 0x20FF) { // Combining Diacritical Marks for Symbols
			runes = append(runes, r)
		}
	}
	return string(runes)
}
