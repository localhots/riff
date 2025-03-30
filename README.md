# Riff

Riff is a logger.

```go
import "github.com/localhots/riff"
import "github.com/localhots/riff/ctx/log"

log.Setup(riff.DefaultConfig())

log.Info(ctx, "Starting task",
	log.Str("device_unique_id", "G4000E-1000-F"),
	log.Int("task_id", 123456),
	log.Str("status", "success"),
	log.Str("template_name", "index.tpl"),
)
```
