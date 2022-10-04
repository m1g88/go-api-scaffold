package logger

import (
	"bytes"
	"fmt"
	"io"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	Zap *zap.Logger
}

func New() logger {
	z, _ := zap.NewProduction()

	return logger{Zap: z}
}

func (l *logger) RecoveryWithZap() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(l.Zap, true)
}

func (l *logger) MiddlewareLogger() gin.HandlerFunc {
	return ginzap.GinzapWithConfig(l.Zap, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			blw := &BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw

			// log request ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			// log trace and span ID
			if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
				fields = append(fields, zap.String("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
				fields = append(fields, zap.String("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
			}

			// log request body
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			fields = append(fields, zap.String("body", string(body)))

			return fields
		}),
	})
}

type BodyLogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	fmt.Println("###", string(b))
	return io.MultiWriter(w.Body, w.ResponseWriter).Write(b)
}
