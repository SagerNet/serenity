# Introduction

serenity uses JSON for configuration files.

### Structure

```json
{
  "log": {},
  "listen": "",
  "tls": {},
  "cache_file": "",
  "outbounds": [],
  "subscriptions": [],
  "templates": [],
  "profiles": [],
  "users": []
}
```

### Fields

#### log

Log configuration, see [Log](https://sing-box.sagernet.org/configuration/log/).

#### listen

==Required==

Listen address.

#### tls

TLS configuration, see [TLS](https://sing-box.sagernet.org/configuration/shared/tls/#inbound).

#### cache_file

Cache file path.

`cache.db` will be used if empty.

#### outbounds

List of [Outbound][outbound], can be referenced in [Profile](./profile).

For chained outbounds, use an array of outbounds as an item, and the first outbound will be the entry.

#### subscriptions

List of [Subscription](./subscription), can be referenced in [Profile](./profile).

#### templates

==Required==

List of [Template](./template), can be referenced in [Profile](./profile).

#### profiles

==Required==

List of [Profile](./profile), can be referenced in [User](./user).

#### users

==Required==

List of [User](./user).

### Check

```bash
serenity check
```

### Format

```bash
serenity format -w -c config.json -D config_directory
```

[outbound]: https://sing-box.sagernet.org/configuration/outbound/