### Structure

```json
{
  "name": "",
  "password": "",
  "profile": [],
  "default_profile": ""
}
```

!!! note ""

    You can ignore the JSON Array [] tag when the content is only one item

### Fields

#### name

HTTP basic authentication username.

#### password

HTTP basic authentication password.

#### profile

Accessible profiles for this user.

List of [Profile](./profile) name.

#### default_profile

Default profile name.

First profile is used by default.
