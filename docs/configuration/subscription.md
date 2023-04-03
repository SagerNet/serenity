# Subscription

A list of external sources to pull server groups from.

### Structure

```json
{
  "subscriptions": [
    {
      "name": "my", // required
      "url": "http://example.org", // required
      "user_agent": "",
      "update_interval": "5m",
      "generate_selector": false,
      "generate_url_test": false
    }
  ]
}
```

### Subscription support

* Clash configuration
* SIP008 delivery document

### Fields

#### name

===Required===

The name of the subscription, will be used in generating group outbounds like selectors.

#### url

===Required===

The URL of the subscription.

#### user_agent

The user agent of the subscription request.

`ClashForAndroid/serenity` is used by default.

#### update_interval

The interval of the subscription update.

5 minutes is used by default.

#### generate_selector

Whether to generate selector outbound for the subscription.

#### generate_url_test

Whether to generate URL test outbound for the subscription.