package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type PrettyJSONEncoder struct {
	zapcore.Encoder
}

// EncodeEntry Override Default Encoding with Prettier Encoding
func (e *PrettyJSONEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	var err error

	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		panic("Failed to encode entry" + err.Error())
	}

	var prettyJSON map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &prettyJSON); err != nil {
		return buf, nil // fallback if fail
	}

	// Prettier raw json
	prettyJson, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		panic("Failed to prettier raw json" + err.Error())
	}

	buf.Reset()
	buf.AppendString(string(prettyJson) + "\n")
	return buf, nil
}

func InitLogger() *zap.Logger {
	var zapLogger *zap.Logger
	var err error

	if Configs.Mode == "prod" || Configs.Mode == "production" {
		zapLogger, err = zap.NewProduction()
		if err != nil {
			panic("Failed to create production logger: " + err.Error())
		}

	} else {
		// Default: development Configs.Mode
		zapLogger = newPrettyLogger()
	}

	zap.ReplaceGlobals(zapLogger)
	zap.L().Info("Logger initialized")
	return zapLogger
}

func UseGinRequestLogging(engine *gin.Engine) {
	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		data := customLoggingFormat(param)

		var jsonData []byte
		var err error
		if Configs.Mode == "prod" || Configs.Mode == "production" {
			jsonData, err = json.Marshal(data)
		} else {
			jsonData, err = json.MarshalIndent(data, "", "  ")
		}

		if err != nil {
			return fmt.Sprintf(`{"level":"error","msg":"failed to marshal log","error":"%v"}`+"\n", err)
		}
		return string(jsonData) + "\n"
	}))

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		zap.L().Info(fmt.Sprintf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers))
	}
}

func newPrettyLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)
	prettyEncoder := &PrettyJSONEncoder{jsonEncoder}

	core := zapcore.NewCore(prettyEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
	return zap.New(core, zap.AddCaller())
}

func customLoggingFormat(param gin.LogFormatterParams) map[string]interface{} {
	return map[string]interface{}{
		"ts":         param.TimeStamp.Format(time.RFC3339Nano),
		"level":      "info",
		"logger":     "zap.L()",
		"msg":        "HTTP request",
		"client_ip":  param.ClientIP,
		"method":     param.Method,
		"path":       param.Path,
		"status":     param.StatusCode,
		"latency_ms": param.Latency.Milliseconds(),
		"user_agent": param.Request.UserAgent(),
		"error":      param.ErrorMessage,
		"proto":      param.Request.Proto,
		"bytes_in":   param.Request.Header.Get("Content-Length"),
		"bytes_out":  param.BodySize,
	}
}
