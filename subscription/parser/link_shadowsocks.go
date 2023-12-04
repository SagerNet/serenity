package parser

import (
	"net/url"
	"strconv"
	"strings"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func ParseShadowsocksLink(link string) (option.Outbound, error) {
	linkURL, err := url.Parse(link)
	if err != nil {
		return option.Outbound{}, err
	}

	if linkURL.User == nil {
		return option.Outbound{}, E.New("missing user info")
	}

	var options option.ShadowsocksOutboundOptions
	options.ServerOptions.Server = linkURL.Host
	options.ServerOptions.ServerPort = portFromString(linkURL.Port())
	if password, _ := linkURL.User.Password(); password != "" {
		options.Method = linkURL.User.Username()
		options.Password = password
	} else {
		userAndPassword, _ := decodeBase64URLSafe(linkURL.User.Username())
		userAndPasswordParts := strings.Split(userAndPassword, ":")
		if len(userAndPasswordParts) != 2 {
			return option.Outbound{}, E.New("bad user info")
		}
		options.Method = userAndPasswordParts[0]
		options.Password = userAndPasswordParts[1]
	}

	plugin := linkURL.Query().Get("plugin")
	options.Plugin = shadowsocksPluginName(plugin)
	options.PluginOptions = shadowsocksPluginOptions(plugin)

	var outbound option.Outbound
	outbound.Type = C.TypeShadowsocks
	outbound.Tag = linkURL.Fragment
	outbound.ShadowsocksOptions = options
	return outbound, nil
}

func portFromString(portString string) uint16 {
	port, _ := strconv.ParseUint(portString, 10, 16)
	return uint16(port)
}

func shadowsocksPluginName(plugin string) string {
	if index := strings.Index(plugin, ";"); index != -1 {
		return plugin[:index]
	}
	return plugin
}

func shadowsocksPluginOptions(plugin string) string {
	if index := strings.Index(plugin, ";"); index != -1 {
		return plugin[index+1:]
	}
	return ""
}
