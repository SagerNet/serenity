package template

import (
	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func (t *Template) renderGeoResources(metadata M.Metadata, options *option.Options) {
	if t.DisableRuleSet || (metadata.Version == nil || metadata.Version.LessThan(semver.ParseVersion("1.8.0-alpha.10"))) {
		var (
			geoipDownloadURL   string
			geositeDownloadURL string
			downloadDetour     string
		)
		if t.EnableJSDelivr {
			geoipDownloadURL = "https://testingcf.jsdelivr.net/gh/SagerNet/sing-geoip@release/geoip-cn.db"
			geositeDownloadURL = "https://testingcf.jsdelivr.net/gh/SagerNet/sing-geosite@release/geosite-cn.db"
			if t.DirectTag != "" {
				downloadDetour = t.DirectTag
			} else {
				downloadDetour = DefaultDirectTag
			}
		} else {
			geoipDownloadURL = "https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip-cn.db"
			geositeDownloadURL = "https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite-cn.db"
		}
		if t.CustomGeoIP == nil {
			options.Route.GeoIP = &option.GeoIPOptions{
				DownloadURL:    geoipDownloadURL,
				DownloadDetour: downloadDetour,
			}
		}
		if t.CustomGeosite == nil {
			options.Route.Geosite = &option.GeositeOptions{
				DownloadURL:    geositeDownloadURL,
				DownloadDetour: downloadDetour,
			}
		}
	} else if len(t.CustomDNSRules) == 0 {
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
		options.Route.RuleSet = []option.RuleSet{
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geoip-cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL:            downloadURL + "SagerNet/sing-geoip" + branchSplit + "rule-set/geoip-cn.srs",
					DownloadDetour: downloadDetour,
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL:            downloadURL + "SagerNet/sing-geosite" + branchSplit + "rule-set/geosite-cn.srs",
					DownloadDetour: downloadDetour,
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-geolocation-!cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL:            downloadURL + "SagerNet/sing-geosite" + branchSplit + "rule-set/geosite-geolocation-!cn.srs",
					DownloadDetour: downloadDetour,
				},
			},
			{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-category-companies@cn",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL:            downloadURL + "SagerNet/sing-geosite" + branchSplit + "rule-set/geosite-category-companies@cn.srs",
					DownloadDetour: downloadDetour,
				},
			},
		}
		if metadata.Platform.IsApple() {
			options.Route.RuleSet = append(options.Route.RuleSet, option.RuleSet{
				Type:   C.RuleSetTypeRemote,
				Tag:    "geosite-apple-update",
				Format: C.RuleSetFormatBinary,
				RemoteOptions: option.RemoteRuleSet{
					URL:            downloadURL + "SagerNet/sing-geosite" + branchSplit + "rule-set/geosite-apple-update.srs",
					DownloadDetour: downloadDetour,
				},
			})
		}
	}
}
