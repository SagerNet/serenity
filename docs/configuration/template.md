### Structure

```json
{
  "name": "",
  "extend": "",
  // Global

  "log": {},
  "domain_strategy": "",
  "domain_strategy_local": "",
  "disable_traffic_bypass": false,
  "disable_rule_set": false,
  "remote_resolve": false,
  // DNS

  "dns": "",
  "dns_local": "",
  "enable_fakeip": false,
  "pre_dns_rules": [],
  "custom_dns_rules": [],
  // Inbound

  "inbounds": [],
  "auto_redirect": false,
  "disable_tun": false,
  "disable_system_proxy": false,
  "custom_tun": {},
  "custom_mixed": {},
  // Outbound

  "extra_groups": [
    {
      "tag": "",
      "type": "",
      "filter": "",
      "exclude": "",
      "custom_selector": {},
      "custom_urltest": {}
    }
  ],
  "generate_global_urltest": false,
  "direct_tag": "",
  "default_tag": "",
  "urltest_tag": "",
  "custom_direct": {},
  "custom_selector": {},
  "custom_urltest": {},
  // Route

  "pre_rules": [],
  "custom_rules": [],
  "enable_jsdelivr": false,
  "custom_geoip": {},
  "custom_geosite": {},
  "custom_rule_set": [],
  "post_rule_set": [],
  // Experimental

  "disable_cache_file": false,
  "disable_clash_mode": false,
  "clash_mode_rule": "",
  "clash_mode_global": "",
  "clash_mode_direct": "",
  "custom_clash_api": {},
  // Debug

  "pprof_listen": "",
  "memory_limit": ""
}
```

### Fields

#### name

==Required==

Profile name.

#### extend

Extend from another profile.

#### log

Log configuration, see [Log](https://sing-box.sagernet.org/configuration/log/).

#### domain_strategy

Global sing-box domain strategy.

One of `prefer_ipv4` `prefer_ipv6` `ipv4_only` `ipv6_only`.

If `*_only` enabled, TUN and DNS will be configured to disable the other network.

Note that if want `prefer_*` to take effect on transparent proxy requests, set `enable_fakeip`.

`ipv4_only` is used by default when `enable_fakeip` disabled,
`prefer_ipv4` is used by default when `enable_fakeip` enabled.

#### domain_strategy_local

Local sing-box domain strategy.

`prefer_ipv4` is used by default.

#### disable_rule_set

Use `geoip` and `geosite` for traffic bypassing instead of rule sets.

#### disable_traffic_bypass

Disable traffic bypass for Chinese DNS queries and connections.

#### remote_resolve

Don't generate `doamin_strategy` options for inbounds.

#### dns

Default DNS server.

`tls://8.8.8.8` is used by default.

#### dns_local

DNS server used for China DNS requests.

`114.114.114.114` is used by default.

#### enable_fakeip

Enable FakeIP.

#### pre_dns_rules

List of [DNS Rule](https://sing-box.sagernet.org/configuration/dns/rule/).

Will be applied before traffic bypassing rules.

#### custom_dns_rules

List of [DNS Rule](https://sing-box.sagernet.org/configuration/dns/rule/).

No default traffic bypassing DNS rules will be generated if not empty.

#### inbounds

List of [Inbound](https://sing-box.sagernet.org/configuration/inbound/).

#### auto_redirect

Generate [auto-redirect](https://sing-box.sagernet.org/configuration/inbound/tun/#auto_redirect) options for android and unknown platforms.

#### disable_tun

Don't generate TUN inbound.

If the target platform can only use TUN for proxy (currently all Apple platforms), this item will not take effect.

#### disable_system_proxy

Don't generate `tun.platform.http_proxy` for known platforms and `set_system_proxy` for unknown platforms.

#### custom_tun

Custom [TUN](https://sing-box.sagernet.org/configuration/inbound/tun/) inbound template.

#### custom_mixed

Custom [Mixed](https://sing-box.sagernet.org/configuration/inbound/mixed/) inbound template.

#### extra_groups

Generate extra outbound groups.

#### extra_groups.tag

==Required==

Tag of the group outbound.

#### extra_groups.type

==Required==

Type of the group outbound.

#### extra_groups.filter

Regexp filter rules, non-matching outbounds will be removed.

#### extra_groups.exclude

Regexp exclude rules, matching outbounds will be removed.

#### extra_groups.custom_selector

Custom [Selector](https://sing-box.sagernet.org/configuration/outbound/selector/) template.

#### extra_groups.custom_urltest

Custom [URLTest](https://sing-box.sagernet.org/configuration/outbound/urltest/) template.

#### generate_global_urltest

Generate a global `URLTest` outbound with all global outbounds.

#### direct_tag

Custom direct outbound tag.

#### default_tag

Custom default outbound tag.

#### urltest_tag

Custom URLTest outbound tag.

#### custom_direct

Custom [Direct](https://sing-box.sagernet.org/configuration/outbound/direct/) outbound template.

#### custom_selector

Custom [Selector](https://sing-box.sagernet.org/configuration/outbound/selector/) outbound template.

#### custom_urltest

Custom [URLTest](https://sing-box.sagernet.org/configuration/outbound/urltest/) outbound template.

#### pre_rules

List of [Rule](https://sing-box.sagernet.org/configuration/route/rule/).

Will be applied before traffic bypassing rules.

#### custom_rules

List of [Rule](https://sing-box.sagernet.org/configuration/route/rule/).

No default traffic bypassing rules will be generated if not empty.

#### enable_jsdelivr

Use jsDelivr CDN and direct outbound for default rule sets or Geo resources.

#### custom_geoip

Custom [GeoIP](https://sing-box.sagernet.org/configuration/route/geoip/) template.

#### custom_geosite

Custom [GeoSite](https://sing-box.sagernet.org/configuration/route/geosite/) template.

#### custom_rule_set

List of [RuleSet](/configuration/shared/rule-set/).

Default rule sets will not be generated if not empty.

#### post_rule_set

List of [RuleSet](/configuration/shared/rule-set/).

Will be applied after default rule sets.

#### disable_cache_file

Don't generate `cache_file` related options.

#### disable_clash_mode

Don't generate `clash_mode` related options.

#### clash_mode_rule

Name of the 'Rule' Clash mode.

`Rule` is used by default.

#### clash_mode_global

Name of the 'Global' Clash mode.

`Global` is used by default.

#### clash_mode_direct

Name of the 'Direct' Clash mode.

`Direct` is used by default.

#### custom_clash_api

Custom [Clash API](https://sing-box.sagernet.org/configuration/experimental/clash-api/) template.

#### pprof_listen

Listen address of the pprof server.

#### memory_limit

Set soft memory limit for sing-box.

`100m` is recommended if memory limit is required.
