# Logging Convention

## In Handler
Use the `ctx.Logger` for request-scoped logs.

```go
// Info | Warn (both same convention)
ctx.Logger.Info("Message")

ctx.Logger.Info(fmt.Sprintf("Message %s, %s", a, b))

ctx.Logger.Info("Message", zap.Any("key", value))
```

## In Service, Repository, or Other Layers
Use the `global logger zap.L()` for system-wide logs.
```go
// Info | Warn (both same convention)
zap.L().Info("Message")

zap.L().Info(fmt.Sprintf("Message %s, %s", a, b))

zap.L().Info("Message", zap.Any("key", value))
```

## Notes
Always prefer structured logging with fields (zap.Any, zap.String, etc.).

Use fmt.Sprintf only when concatenation is necessary.

Keep log messages short, consistent, and descriptive.