package process

import (
    "context"
    "google.golang.org/grpc"
    runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// MyRuntime implements the CRI RuntimeService
type MyRuntime struct {
    runtimeapi.UnimplementedRuntimeServiceServer
}

// Implement Version() to respond to crictl version/info
func (s *MyRuntime) Version(ctx context.Context, req *runtimeapi.VersionRequest) (*runtimeapi.VersionResponse, error) {
    return &runtimeapi.VersionResponse{
        Version:           "0.1.0",
        RuntimeName:       "my-runtime",
        RuntimeVersion:    "0.1.0",
        RuntimeApiVersion: "v1",
    }, nil
}

// RegisterMyRuntime registers the runtime implementation with a gRPC server
func RegisterMyRuntime(grpcServer *grpc.Server) {
    runtimeapi.RegisterRuntimeServiceServer(grpcServer, &MyRuntime{})
}
