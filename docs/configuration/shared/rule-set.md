# RuleSet

RuleSet generate configuration.

### Structure

=== "Default"

    ```json
    {
      "type": "", // optional
      
      ... // Default Fields
    }
    ```

=== "GitHub"

    ```json
    {
      "type": "github",
      "repository": "",
      "path": "",
      "rule-set": []
    }
    ```

=== "Example"

     ```json
     {
       "type": "github",
       "repository": "SagerNet/sing-geosite",
       "path": "rule-set",
       "prefix": "geosite-",
       "rule-set": [
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
      "path": "sing/geo/geosite",
      "prefix": "geosite-",
      "rule-set": [
        "apple",
        "microsoft",
        "openai"
      ]
    }
    ```

### Default Fields

See [RuleSet](https://sing-box.sagernet.org/configuration/rule-set/).

### GitHub Fields

#### repository

GitHub repository, `SagerNet/sing-<geoip/geosite>` or `MetaCubeX/meta-rules-dat`.

#### path

Branch and directory path, `rule-set` or `sing/geo/<geoip/geosite>`.

#### prefix

File prefix, `geoip-` or `geosite-`.

#### rule-set

RuleSet name list.
