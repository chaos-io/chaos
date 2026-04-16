# Logs

`logs` 包基于 `zap` 提供统一日志能力，支持控制台/文件输出和运行时动态调整日志级别。

## 动态日志级别

当同时配置 `levelPattern` 和 `levelPort` 时，logger 会启动一个 HTTP 服务，暴露 `zap.AtomicLevel.ServeHTTP`：

- `GET <levelPattern>`：获取当前日志级别
- `PUT <levelPattern>`：动态修改日志级别

支持两种 `PUT` 请求体：

- `application/json`，例如 `{"level":"debug"}`
- `application/x-www-form-urlencoded`，例如 `level=debug`

## 配置示例

```yaml
logs:
  level: info
  encode: console
  output: console
  levelPattern: /log/level
  levelPort: 22001
```

## 使用方式

应用启动时先加载配置，然后刷新日志默认服务：

```go
_ = config.InitDefault(config.WithWatcherDisabled())
_ = config.LoadPath("./configs")
_ = logs.ReloadDefaultServiceFromConfig()
```

完成后可通过 HTTP 调整级别：

```bash
# 查看当前级别
curl http://127.0.0.1:22001/log/level

# 设置为 debug（JSON）
curl -X PUT \
  -H "Content-Type: application/json" \
  -d '{"level":"debug"}' \
  http://127.0.0.1:22001/log/level

# 设置为 warn（表单）
curl -X PUT \
  -d "level=warn" \
  http://127.0.0.1:22001/log/level
```
