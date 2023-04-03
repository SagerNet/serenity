# Introduction

serenity uses JSON for configuration files.

### Structure

```json
{
  "log": {},
  "listen": "[::]:443",
  "tls": {},
  "subscriptions": [],
  "outbounds": [],
  "profiles": [],
  "default_profile": ""
}
```

#### log

Log configuration, see [Log](https://sing-box.sagernet.org/configuration/log/).

#### listen

Listen address.

#### tls

TLS configuration, see [TLS](https://sing-box.sagernet.org/configuration/shared/tls/#inbound).

#### subscriptions

List of [Subscription](./subscription).

#### outbounds

List of [Outbound](https://sing-box.sagernet.org/configuration/outbound/).

#### profiles

List of [Profile](./profile).

#### default_profile

Default profile.

First profile will be used if empty.