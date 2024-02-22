package template

import (
	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/json/badjson"
)

func (t *Template) renderExperimental(metadata M.Metadata, options *option.Options, profileName string) error {
	if t.DisableCacheFile && t.DisableClashMode && t.CustomClashAPI == nil {
		return nil
	}
	options.Experimental = &option.ExperimentalOptions{}
	disable18Features := metadata.Version != nil && metadata.Version.LessThan(semver.ParseVersion("1.8.0-alpha.10"))
	if !t.DisableCacheFile {
		if disable18Features {
			//nolint:staticcheck
			//goland:noinspection GoDeprecation
			options.Experimental.ClashAPI = &option.ClashAPIOptions{
				CacheID:       profileName,
				StoreMode:     true,
				StoreSelected: true,
				StoreFakeIP:   t.EnableFakeIP,
			}
		} else {
			options.Experimental.CacheFile = &option.CacheFileOptions{
				Enabled:     true,
				CacheID:     profileName,
				StoreFakeIP: t.EnableFakeIP,
			}
			if !t.DisableDNSLeak && (metadata.Version != nil && metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.9.0-alpha.8"))) {
				options.Experimental.CacheFile.StoreRDRC = true
			}
		}
	}

	if t.CustomClashAPI != nil {
		newClashOptions, err := badjson.MergeFromDestination(options.Experimental.ClashAPI, t.CustomClashAPI.Message)
		if err != nil {
			return err
		}
		options.Experimental.ClashAPI = newClashOptions
	} else if options.Experimental.ClashAPI == nil {
		options.Experimental.ClashAPI = &option.ClashAPIOptions{}
	}

	if !t.DisableExternalController && options.Experimental.ClashAPI.ExternalController == "" {
		options.Experimental.ClashAPI.ExternalController = "127.0.0.1:9090"
	}

	if !t.DisableClashMode {
		if !t.DisableDNSLeak && (metadata.Version != nil && metadata.Version.GreaterThanOrEqual(semver.ParseVersion("1.9.0-alpha.1"))) {
			clashModeLeak := t.ClashModeLeak
			if clashModeLeak == "" {
				clashModeLeak = "Leak"
			}
			options.Experimental.ClashAPI.DefaultMode = clashModeLeak
		} else {
			options.Experimental.ClashAPI.DefaultMode = t.ClashModeRule
		}
	}
	if t.PProfListen != "" {
		if options.Experimental.Debug == nil {
			options.Experimental.Debug = &option.DebugOptions{}
		}
		options.Experimental.Debug.Listen = t.PProfListen
	}
	if t.MemoryLimit > 0 && !metadata.Platform.IsNetworkExtensionMemoryLimited() {
		if options.Experimental.Debug == nil {
			options.Experimental.Debug = &option.DebugOptions{}
		}
		options.Experimental.Debug.MemoryLimit = t.MemoryLimit
		options.Experimental.Debug.OOMKiller = common.Ptr(true)
	}
	return nil
}
