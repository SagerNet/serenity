package metadata

import (
	"strings"

	"github.com/sagernet/serenity/common/semver"
	E "github.com/sagernet/sing/common/exceptions"
)

type Platform string

const (
	PlatformUnknown   Platform = ""
	PlatformAndroid   Platform = "android"
	PlatformiOS       Platform = "ios"
	PlatformMacOS     Platform = "macos"
	PlatformAppleTVOS Platform = "tvos"
)

func ParsePlatform(name string) (Platform, error) {
	switch strings.ToLower(name) {
	case "android":
		return PlatformAndroid, nil
	case "ios":
		return PlatformiOS, nil
	case "macos":
		return PlatformMacOS, nil
	case "tvos":
		return PlatformAppleTVOS, nil
	default:
		return PlatformUnknown, E.New("unknown platform: ", name)
	}
}

func (m Platform) IsApple() bool {
	switch m {
	case PlatformiOS, PlatformMacOS, PlatformAppleTVOS:
		return true
	default:
		return false
	}
}

func (m Platform) IsNetworkExtensionMemoryLimited() bool {
	switch m {
	case PlatformiOS, PlatformAppleTVOS:
		return true
	default:
		return false
	}
}

func (m Platform) TunOnly() bool {
	return m.IsApple()
}

func (m Platform) String() string {
	return string(m)
}

type Metadata struct {
	UserAgent string
	Platform  Platform
	Version   *semver.Version
}

func Detect(userAgent string) Metadata {
	var metadata Metadata
	metadata.UserAgent = userAgent
	if strings.HasPrefix(userAgent, "SFA") {
		metadata.Platform = PlatformAndroid
	} else if strings.HasPrefix(userAgent, "SFI") {
		metadata.Platform = PlatformiOS
	} else if strings.HasPrefix(userAgent, "SFM") {
		metadata.Platform = PlatformMacOS
	} else if strings.HasPrefix(userAgent, "SFT") {
		metadata.Platform = PlatformAppleTVOS
	}
	var versionName string
	if strings.Contains(userAgent, "sing-box ") {
		versionName = strings.Split(userAgent, "sing-box ")[1]
		if strings.Contains(versionName, ";") {
			versionName = strings.Split(versionName, ";")[0]
		} else {
			versionName = strings.Split(versionName, ")")[0]
		}
	}
	if semver.IsValid(versionName) {
		version := semver.ParseVersion(versionName)
		metadata.Version = &version
	}
	return metadata
}
