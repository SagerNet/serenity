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
      "filter_type": [],
      "exclude_type": [],
      "invert": false,
      "remove": false,
      "rename": {},
      "remove_emoji": false,
      "rewrite_multiplex": {}
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

`serenity/$version (sing-box $sing-box-version; Clash compatible)` is used by default.

#### process

!!! note ""

    You can ignore the JSON Array [] tag when the content is only one item

Process rules.

#### process.filter

Regexp filter rules, match outbound tag name.

#### process.exclude

Regexp exclude rules, match outbound tag name.

#### process.filter_type

Filter rules, match outbound type.

#### process.exclude_type

Exclude rules, match outbound type.

#### process.invert

Invert filter results.

#### process.remove

Remove outbounds that match the rules.

#### process.rename

Regexp rename rules, matching outbounds will be renamed.

#### process.remove_emoji

Remove emojis in outbound tags.

#### process.rewrite_multiplex

Rewrite [Multiplex](https://sing-box.sagernet.org/configuration/shared/multiplex) options.

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
