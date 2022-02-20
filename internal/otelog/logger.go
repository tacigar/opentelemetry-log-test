package otelog

import (
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	projectId = "dev"
)

type LogContent struct {
	Message string
	Span    trace.Span
}

type Logger interface {
	Info(content LogContent)
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: logger}, nil
}

func (l *ZapLogger) Info(content LogContent) {
	l.logger.Info(content.Message,
		zap.String(
			"logging.googleapis.com/trace",
			fmt.Sprintf(
				"projects/%s/traces/%s",
				projectId,
				content.Span.SpanContext().TraceID().String(),
			),
		),
		zap.String("logging.googleapis.com/spanId",
			content.Span.SpanContext().SpanID().String(),
		),
		zap.Bool("logging.googleapis.com/trace_sampled",
			content.Span.SpanContext().IsSampled(),
		),
	)
}
