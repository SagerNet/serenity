---
icon: material/package
---

# Package Manager

## :material-download-box: Manual Installation

=== ":material-debian: Debian / DEB"

    ```bash
    bash <(curl -fsSL https://sing-box.app/serenity/deb-install.sh)
    ```

=== ":material-redhat: Redhat / RPM"

    ```bash
    bash <(curl -fsSL https://sing-box.app/serenity/rpm-install.sh)
    ```

=== ":simple-archlinux: Archlinux / PKG"

    ```bash
    bash <(curl -fsSL https://sing-box.app/serenity/arch-install.sh)
    ```

## :material-book-multiple: Service Management

For Linux systems with [systemd][systemd], usually the installation already includes a serenity service,
you can manage the service using the following command:

| Operation | Command                                       |
|-----------|-----------------------------------------------|
| Enable    | `sudo systemctl enable serenity`              |
| Disable   | `sudo systemctl disable serenity`             |
| Start     | `sudo systemctl start serenity`               |
| Stop      | `sudo systemctl stop serenity`                |
| Kill      | `sudo systemctl kill serenity`                |
| Restart   | `sudo systemctl restart serenity`             |
| Logs      | `sudo journalctl -u serenity --output cat -e` |
| New Logs  | `sudo journalctl -u serenity --output cat -f` |

[systemd]: https://systemd.io/