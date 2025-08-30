# RuleSet

RuleSet generate configuration.

### Structure

=== "Original"

    ```json
    {
      "type": "remote", // or local
      
      ... // Original Fields
    }
    ```

=== "GitHub"

    ```json
    {
      "type": "github",
      "repository": "",
      "path": "",
      "rule_set": []
    }
    ```

=== "Example"

     ```json
     {
       "type": "github",
       "repository": "SagerNet/sing-geosite",
       "path": "rule-set/geosite-",
       "prefix": "geosite-",
       "rule_set": [
         "apple",
         "microsoft",
         "openai"
       ]
     }
     ```

=== "Example (Clash.Meta repository)"

    ```json
    {
      "type": "github",
      "repository": "MetaCubeX/meta-rules-dat",
      "path": "sing/geo/geosite/",
      "prefix": "geosite-",
      "rule_set": [
        "apple",
        "microsoft",
        "openai"
      ]
    }
    ```

### Original Fields

See [RuleSet](https://sing-box.sagernet.org/configuration/rule-set/).

### GitHub Fields

#### repository

GitHub repository, `SagerNet/sing-<geoip/geosite>` or `MetaCubeX/meta-rules-dat`.

#### path

Branch and directory path, `rule-set` or `sing/geo/<geoip/geosite>`.

#### prefix

File prefix, `geoip-` or `geosite-`.

#### rule_set

RuleSet name list.
