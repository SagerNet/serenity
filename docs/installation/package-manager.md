---
icon: material/package
---

# Package Manager

## :material-tram: Repository Installation

=== ":material-debian: Debian / APT"

    ```bash
    sudo curl -fsSL https://deb.sagernet.org/gpg.key -o /etc/apt/keyrings/sagernet.asc
    sudo chmod a+r /etc/apt/keyrings/sagernet.asc
    echo "deb [arch=`dpkg --print-architecture` signed-by=/etc/apt/keyrings/sagernet.asc] https://deb.sagernet.org/ * *" | \
      sudo tee /etc/apt/sources.list.d/sagernet.list > /dev/null
    sudo apt-get update
    sudo apt-get install serenity
    ```

=== ":material-redhat: Redhat / DNF"

    ```bash
    sudo dnf -y install dnf-plugins-core
    sudo dnf config-manager --add-repo https://sing-box.app/rpm.repo
    sudo dnf install serenity
    ```

=== ":material-redhat: CentOS / YUM"

    ```bash
    sudo yum install -y yum-utils
    sudo yum-config-manager --add-repo https://sing-box.app/rpm.repo
    sudo yum install serenity
    ```

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
