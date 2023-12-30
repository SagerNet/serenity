### Structure

```json
{
  "name": "",
  "url": "",
  "user_agent": "",
  "process": [
    {
      "filter": [],
      "exclude": [],
      "filter_outbound_type": [],
      "exclude_outbound_type": [],
      "rename": {},
      "remove_emoji": false
    }
  ],
  "deduplication": false,
  "update_interval": "5m",
  "generate_selector": false,
  "generate_urltest": false,
  "urltest_suffix": "",
  "custom_selector": {},
  "custom_urltest": {}
}
```

### Fields

#### name

==Required==

Name of the subscription, will be used in group tags.

#### url

==Required==

Subscription URL.

#### user_agent

User-Agent in HTTP request.

`ClashForAndroid/serenity` is used by default.

#### process

!!! note ""

    You can ignore the JSON Array [] tag when the content is only one item

Process rules.

#### process.filter

Regexp filter rules, non-matching outbounds will be removed.

#### process.exclude

Regexp exclude rules, matching outbounds will be removed.

#### process.filter_outbound_type

Outbound type filter rules, non-matching outbounds will be removed.

#### process.exclude_outbound_type

Outbound type exclude rules, matching outbounds will be removed.

#### process.rename

Regexp rename rules, matching outbounds will be renamed.

#### process.remove_emoji

Remove emojis in outbound tags.

#### deduplication

Remove outbounds with duplicate server destinations (Domain will be resolved to compare).

#### update_interval

Subscription update interval.

`1h` is used by default.

#### generate_selector

Generate a global `Selector` outbound for the subscription.

If both `generate_selector` and `generate_urltest` are disabled, subscription outbounds will be added to global groups.

#### generate_urltest

Generate a global `URLTest` outbound for the subscription.

If both `generate_selector` and `generate_urltest` are disabled, subscription outbounds will be added to global groups.

#### urltest_suffix

Tag suffix of generated `URLTest` outbound.

` - URLTest` is used by default.

#### custom_selector

Custom [Selector](https://sing-box.sagernet.org/configuration/outbound/selector/) template.

#### custom_urltest

Custom [URLTest](https://sing-box.sagernet.org/configuration/outbound/urltest/) template.
