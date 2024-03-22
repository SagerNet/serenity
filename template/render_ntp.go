package template

import (
	"time"

	"github.com/sagernet/sing-box/option"
)

func (t *Template) RenderNTP(options *option.Options) error {
	if t.EnableNTP {
		options.NTP = &option.NTPOptions{
			Enabled:    true,
			ServerPort: 123,
			Server:     "time.apple.com",
			Interval:   option.Duration(time.Minute * 30),
		}

		if t.CustomNTP != nil {
			if t.CustomNTP.Server != "" {
				options.NTP.Server = t.CustomNTP.Server
			}

			if t.CustomNTP.ServerPort != 0 {
				options.NTP.ServerPort = t.CustomNTP.ServerPort
			}

			if t.CustomNTP.Interval != 0 {
				options.NTP.Interval = t.CustomNTP.Interval
			}

			if t.CustomNTP.DialerOptions != (option.DialerOptions{}) {
				options.NTP.DialerOptions = t.CustomNTP.DialerOptions
			}
		}

	}
	return nil
}
