package datadog

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// LogrusDDTraceContextInjector adds dd.trace_id, dd.span_id as logrus entry fields and puts new log entry into the context. Other interceptors can pick the log entry with ctxlogrus.Extract and add extra data. This interceptor must be used after grpctrace.UnaryServerInterceptor (from "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc") and before grpc_logrus.UnaryServerInterceptor (from github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus)
//
// Example:
//	ddtracer.Start()
//	defer ddtracer.Flush()
//	l := logrus.New()
//	l.Out = os.Stderr
//	logger := logrus.NewEntry(l)
//	server := grpctest.NewServer(grpc_middleware.WithUnaryServerChain(
// 		grpctrace.UnaryServerInterceptor(),
// 		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
// 		midleware.LogrusDDTraceContextInjector(logger),
// 		grpc_logrus.UnaryServerInterceptor(logger),
// 	))
func LogrusDDTraceContextInjector(entry *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if span, exists := tracer.SpanFromContext(ctx); exists && span.Context().TraceID() != 0 {
			callLog := entry.WithFields(
				logrus.Fields{
					"dd.trace_id": span.Context().TraceID(),
					"dd.span_id":  span.Context().SpanID(),
				})
			callLog = callLog.WithFields(ctxlogrus.Extract(ctx).Data)
			ctx = ctxlogrus.ToContext(ctx, callLog)
		}

		return handler(ctx, req)
	}
}
