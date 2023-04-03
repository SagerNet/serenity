package libsubscription

import "github.com/sagernet/sing-box/option"

type Server struct {
	Name      string            `json:"name"`
	Outbounds []option.Outbound `json:"outbounds"`
}
