### Structure

```json
{
  "name": "",
  "template": "",
  "template_for_platform": {},
  "template_for_user_agent": {},
  "outbound": [],
  "subscription": []
}
```

### Fields

#### name

==Required==

Profile name.

#### template

Default template name.

A empty template is used by default.

#### template_for_platform

Custom template for different graphical client.

The key is one of `android`, `ios`, `macos`, `tvos`.

The Value is the template name.

#### template_for_user_agent

Custom template for different user agent.

The key is a regular expression matching the user agent.

The value is the template name.

#### outbound

Included outbounds.

#### subscription

Included subscriptions.
