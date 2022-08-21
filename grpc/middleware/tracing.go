package middleware

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// UnaryClientTracing 客户端一元链路跟踪
func UnaryClientTracing() grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor()
}

// StreamClientTracing 客户端流链路跟踪
func StreamClientTracing() grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor()
}

// UnaryServerTracing 服务端一元链路跟踪
func UnaryServerTracing() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor()
}

// StreamServerTracing 服务端流链路跟踪
func StreamServerTracing() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor()
}
