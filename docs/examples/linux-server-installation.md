#### Requirements

* Linux & Systemd
* Git
* C compiler environment

#### Install

```shell
git clone -b main https://github.com/SagerNet/serenity
cd serenity
./release/local/install_go.sh # skip if you have golang already installed
./release/local/install.sh
```

Edit configuration file in `/usr/local/etc/serenity/config.json`

```shell
./release/local/enable.sh
```

#### Update

```shell
./release/local/update.sh
```

#### Other commands

| Operation | Command                                       |
|-----------|-----------------------------------------------|
| Start     | `sudo systemctl start serenity`               |
| Stop      | `sudo systemctl stop serenity`                |
| Kill      | `sudo systemctl kill serenity`                |
| Restart   | `sudo systemctl restart serenity`             |
| Logs      | `sudo journalctl -u serenity --output cat -e` |
| New Logs  | `sudo journalctl -u serenity --output cat -f` |
| Uninstall | `./release/local/uninstall.sh`                |