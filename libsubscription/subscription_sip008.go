package libsubscription

import (
	"github.com/sagernet/sing-box/common/json"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

type ShadowsocksDocument struct {
	Version int                         `json:"version"`
	Servers []ShadowsocksServerDocument `json:"servers"`
}

type ShadowsocksServerDocument struct {
	ID         string `json:"id"`
	Remarks    string `json:"remarks"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	Plugin     string `json:"plugin"`
	PluginOpts string `json:"plugin_opts"`
}

func ParseSIP008Subscription(content string) ([]Server, error) {
	var document ShadowsocksDocument
	err := json.Unmarshal([]byte(content), &document)
	if err != nil {
		return nil, E.Cause(err, "parse SIP008 document")
	}

	var servers []Server
	for _, server := range document.Servers {
		servers = append(servers, Server{
			Name: server.Remarks,
			Outbounds: []option.Outbound{
				{
					Type: C.TypeShadowsocks,
					Tag:  server.Remarks,
					ShadowsocksOptions: option.ShadowsocksOutboundOptions{
						ServerOptions: option.ServerOptions{
							Server:     server.Server,
							ServerPort: uint16(server.ServerPort),
						},
						Password:      server.Password,
						Method:        server.Method,
						Plugin:        server.Plugin,
						PluginOptions: server.PluginOpts,
					},
				},
			},
		})
	}
	return servers, nil
}
