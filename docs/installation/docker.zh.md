---
icon: material/docker
---

# Docker

## :material-console: 命令

```bash
docker run -d \
  -v /etc/serenity:/etc/serenity/ \
  --name=serenity \
  --restart=always \
  ghcr.io/sagernet/serenity \
  -D /var/lib/serenity \
  -C /etc/serenity/ run
```

## :material-box-shadow: Compose

```yaml
version: "3.8"
services:
  serenity:
    image: ghcr.io/sagernet/serenity
    container_name: serenity
    restart: always
    volumes:
      - /etc/serenity:/etc/serenity/
    command: -D /var/lib/serenity -C /etc/serenity/ run
```
