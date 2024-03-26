---
icon: material/package
---

# 包管理器

## :material-tram: 仓库安装

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

## :material-download-box: 手动安装

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

## :material-book-multiple: 服务管理

对于带有 [systemd][systemd] 的 Linux 系统，通常安装已经包含 serenity 服务，
您可以使用以下命令管理服务：

| 行动   | 命令                                            |
|------|-----------------------------------------------|
| 启用   | `sudo systemctl enable serenity`              |
| 禁用   | `sudo systemctl disable serenity`             |
| 启动   | `sudo systemctl start serenity`               |
| 停止   | `sudo systemctl stop serenity`                |
| 强行停止 | `sudo systemctl kill serenity`                |
| 重新启动 | `sudo systemctl restart serenity`             |
| 查看日志 | `sudo journalctl -u serenity --output cat -e` |
| 实时日志 | `sudo journalctl -u serenity --output cat -f` |

[systemd]: https://systemd.io/