package serenity

import (
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
)

func filterH2Mux(outbound option.Outbound) bool {
	return !common.Contains([]string{
		common.PtrValueOrDefault(outbound.ShadowsocksOptions.Multiplex).Protocol,
		common.PtrValueOrDefault(outbound.VMessOptions.Multiplex).Protocol,
		common.PtrValueOrDefault(outbound.TrojanOptions.Multiplex).Protocol,
	}, "h2mux")
}

func filterMuxPadding(outbound option.Outbound) bool {
	return !common.Contains([]bool{
		common.PtrValueOrDefault(outbound.ShadowsocksOptions.Multiplex).Padding,
		common.PtrValueOrDefault(outbound.VMessOptions.Multiplex).Padding,
		common.PtrValueOrDefault(outbound.TrojanOptions.Multiplex).Padding,
	}, true)
}
