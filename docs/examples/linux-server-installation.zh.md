#### 依赖

* Linux & Systemd
* Git
* C 编译器环境

#### 安装

```shell
git clone -b main https://github.com/SagerNet/serenity
cd serenity
./release/local/install_go.sh # 如果已安装 golang 则跳过
./release/local/install.sh
```

编辑配置文件 `/usr/local/etc/serenity/config.json`

```shell
./release/local/enable.sh
```

#### 更新

```shell
./release/local/update.sh
```

#### 其他命令

| 操作   | 命令                                            |
|------|-----------------------------------------------|
| 启动   | `sudo systemctl start serenity`               |
| 停止   | `sudo systemctl stop serenity`                |
| 强制停止 | `sudo systemctl kill serenity`                |
| 重启   | `sudo systemctl restart serenity`             |
| 查看日志 | `sudo journalctl -u serenity --output cat -e` |
| 实时日志 | `sudo journalctl -u serenity --output cat -f` |
| 卸载   | `./release/local/uninstall.sh`                |