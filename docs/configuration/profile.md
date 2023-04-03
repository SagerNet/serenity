# Profile

### Structure

```json
{
  "profiles": [
    {
      "name": "default", // required
      "template": "default",
      "config": [
        "/etc/serenity/config.d/config.json"
      ],
      "config_directory": [
        "/etc/serenity/config.d"
      ],
      "group_tag": [
        "select",
        "url-test"
      ],
      "filter_subscription": [
        "none"
      ],
      "filter_outbound": [
        "none"
      ],
      "authorization": {
        "username": "hello",
        "password": "world"
      }
    }
  ]
}
```