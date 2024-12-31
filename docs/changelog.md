---
icon: material/alert-decagram
---

#### 1.1.0-beta.3

* Rename `template.servers` to `template.dns_servers`
* Set `tun.route_exclude_address_set` for traffic bypassing
* Fixes and improvements

#### 1.1.0-alpha.3

* If `extra_groups.include/exclude` is not set,
  subscription groups instead of subscription outbounds will be used as group outbounds
* Add `template.custom_fakeip`
* Fixes and improvements

#### 1.1.0-alpha.1

* Add support for rule actions

#### 1.0.0-beta.19

* Fixes and improvements

#### 1.0.0-beta.17

* Replace `template.extra_groups.<per_subscription/exclude_outbounds>` with `target`

#### 1.0.0-beta.16

* Add `export <profile>` command to export configuration without running the server
* Add `template.extra_groups.exclude_outbounds`
* Add `template.extra_groups.<per_subscription/tag_per_subscription>`

#### 1.0.0-beta.15

* Add support for inline rule-sets **1**

**1**:

Will be merged into route and DNS rules in older versions.

#### 1.0.0-beta.14

* Rename `template.dns_default` to `template.dns`
* Add `template.domain_strategy_local`

#### 1.0.0-beta.13

* Add `template.auto_redirect`

#### 1.0.0-beta.12

* Fixes and improvements

#### 1.0.0-beta.11

* Add `template.extend`
* Add independent rule-set configuration **1**

**1**:

With the new `type=github`, you can batch generate rule-sets based on GitHub files.
See [Rule-Set](/configuration/shared/rule-set/).

#### 1.0.0-beta.10

* Add `template.inbounds`

#### 1.0.0-beta.9

* Add `template.log`
* Fixes and improvements

#### 1.0.0-beta.8

* Add `subscription.process.rewrite_multiplex`
* Rename `subscription.process.[filter/exclude]_outbound_type` to `subscription.process.[filter/exclude]_type`
* Rewrite `subscription.process`
* Fixes and improvements

##### 2023/12/12

No changelog before.
