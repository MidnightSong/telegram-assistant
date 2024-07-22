package utils

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

const DefaultTimeLayout = "2006-01-02 15:04:05"
const TraceID ContextKey = "TraceID"
const Debug ContextKey = "Debug"

type ContextKey string

func init() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(DefaultTimeLayout)

	out2Console := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), zapcore.AddSync(os.Stdout), zapcore.Level(zap.InfoLevel))
	Logger = zap.New(zapcore.NewTee(out2Console), zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
}

// func GenTraceContext() context.Context {
// 	ctx := context.WithValue(context.Background(), TraceID, Hex(GetUUID()[:6]))
// 	if GetConfig().IsDebug() {
// 		return context.WithValue(ctx, Debug, true)
// 	}
// 	return ctx
// }

func fmtCtxMsg(ctx context.Context, msg string) string {
	if ctx != nil {
		if id, _ := ctx.Value(TraceID).(string); id != "" {
			msg = id + "\t" + msg
		}
	}
	return msg
}

func LogDebug(ctx context.Context, msg string, fields ...zap.Field) {
	msg = fmtCtxMsg(ctx, msg)
	Logger.Debug(msg, fields...)
}

func LogInfo(ctx context.Context, msg string, fields ...zap.Field) {
	msg = fmtCtxMsg(ctx, msg)

	Logger.Info(msg, fields...)
}

func LogWarn(ctx context.Context, msg string, fields ...zap.Field) {
	msg = fmtCtxMsg(ctx, msg)
	Logger.Warn(msg, fields...)
}

func LogError(ctx context.Context, msg string, fields ...zap.Field) {
	msg = fmtCtxMsg(ctx, msg)
	Logger.Error(msg, fields...)
}

type RestyLogger struct {
	Proxy string
	Ctx   context.Context
}

func (l *RestyLogger) Errorf(format string, v ...any) {
	LogError(l.Ctx, l.format(format, v...), zap.String("proxy", l.Proxy))
}

func (l *RestyLogger) Warnf(format string, v ...any) {
	LogWarn(l.Ctx, l.format(format, v...), zap.String("proxy", l.Proxy))
}

func (l *RestyLogger) Debugf(format string, v ...any) {
	LogDebug(l.Ctx, l.format(format, v...), zap.String("proxy", l.Proxy))
}

func (l *RestyLogger) format(format string, v ...any) string {
	if len(v) > 0 {
		format = fmt.Sprintf("RESTY LOG"+format, v...)
	}

	str, err := strconv.Unquote(strings.Replace(strconv.Quote(format), `\\u`, `\u`, -1))
	if err != nil {
		LogError(l.Ctx, err.Error())
	}

	return str
}
