package template

import (
	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/constant"
	"github.com/sagernet/serenity/option"
	C "github.com/sagernet/sing-box/constant"
	boxOption "github.com/sagernet/sing-box/option"
)

func (t *Template) renderGeoResources(metadata M.Metadata, options *boxOption.Options) {
	if len(t.CustomRuleSet) == 0 {
		var (
			downloadURL    string
			downloadDetour string
			branchSplit    string
		)
		if t.EnableJSDelivr {
			downloadURL = "https://testingcf.jsdelivr.net/gh/"
			if t.DirectTag != "" {
				downloadDetour = t.DirectTag
			} else {
				downloadDetour = DefaultDirectTag
			}
			branchSplit = "@"
		} else {
			downloadURL = "https://raw.githubusercontent.com/"
			branchSplit = "/"
		}
		options.Route.RuleSet = []boxOption.RuleSet{
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geoip-cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: boxOption.RemoteRuleSet{
					URL:            downloadURL + "SagerNet/sing-geoip" + branchSplit + "rule-set/geoip-cn.srs",
					DownloadDetour: downloadDetour,
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-geolocation-cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: boxOption.RemoteRuleSet{
					URL:            downloadURL + "SagerNet/sing-geosite" + branchSplit + "rule-set/geosite-geolocation-cn.srs",
					DownloadDetour: downloadDetour,
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-geolocation-!cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: boxOption.RemoteRuleSet{
					URL:            downloadURL + "SagerNet/sing-geosite" + branchSplit + "rule-set/geosite-geolocation-!cn.srs",
					DownloadDetour: downloadDetour,
				},
			},
		}
	}
	options.Route.RuleSet = append(options.Route.RuleSet, t.renderRuleSet(t.PostRuleSet)...)
}

func (t *Template) renderRuleSet(ruleSets []option.RuleSet) []boxOption.RuleSet {
	var result []boxOption.RuleSet
	for _, ruleSet := range ruleSets {
		if ruleSet.Type == constant.RuleSetTypeGitHub {
			var (
				downloadURL    string
				downloadDetour string
				branchSplit    string
			)
			if t.EnableJSDelivr {
				downloadURL = "https://testingcf.jsdelivr.net/gh/"
				if t.DirectTag != "" {
					downloadDetour = t.DirectTag
				} else {
					downloadDetour = DefaultDirectTag
				}
				branchSplit = "@"
			} else {
				downloadURL = "https://raw.githubusercontent.com/"
				branchSplit = "/"
			}

			for _, code := range ruleSet.GitHubOptions.RuleSet {
				result = append(result, boxOption.RuleSet{
					Type:   C.RuleSetTypeRemote,
					Tag:    ruleSet.GitHubOptions.Prefix + code,
					Format: C.RuleSetFormatBinary,
					RemoteOptions: boxOption.RemoteRuleSet{
						URL: downloadURL +
							ruleSet.GitHubOptions.Repository +
							branchSplit +
							ruleSet.GitHubOptions.Path +
							"/" +
							ruleSet.GitHubOptions.Prefix +
							code + ".srs",
						DownloadDetour: downloadDetour,
					},
				})
			}
		} else {
			result = append(result, ruleSet.DefaultOptions)
		}
	}
	return result
}
