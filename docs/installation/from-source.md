# Install from source

serenity requires Golang **1.20** or a higher version.

```bash
go install -v -tags with_acme github.com/sagernet/serenity/cmd/serenity@latest
```

The binary is built under $GOPATH/bin

It is also recommended to use systemd to manage serenity service,
see [Linux server installation example](/examples/linux-server-installation).
