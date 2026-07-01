# GoSteamRestarter

跨平台 Steam 重启工具，支持 CLI 和桌面 GUI 两种模式。

## 支持平台

- Windows
- Linux
- macOS

## 构建

```bash
# 构建 CLI 版本
go build ./cmd/cli

# 构建桌面 GUI 版本
go build ./cmd/desktop
```

## 运行

```bash
# 运行 CLI 版本
./cli

# 运行桌面 GUI 版本
./desktop
```

## 测试

```bash
go test ./...
```
