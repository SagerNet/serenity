# 从源代码安装

serenity 需要 Golang **1.20** 或更高版本。

```bash
go install -v -tags with_acme github.com/sagernet/serenity/cmd/serenity@latest
```

二进制文件将被构建在 `$GOPATH/bin` 下。

同时推荐使用 systemd 来管理 serenity 服务器实例。
参阅 [Linux 服务器安装示例](/examples/linux-server-installation)。
