package process

import (
    "context"
    "google.golang.org/grpc"
    runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// MyRuntime implements the CRI RuntimeService
type MyRuntime struct {
    runtimeapi.UnimplementedRuntimeServiceServer
    runtimeapi.UnimplementedImageServiceServer
}

// Implement Version() to respond to crictl version/info
func (s *MyRuntime) Version(ctx context.Context, req *runtimeapi.VersionRequest) (*runtimeapi.VersionResponse, error) {
    return &runtimeapi.VersionResponse{
        Version:           "0.1.0",
        RuntimeName:       "gontainers",
        RuntimeVersion:    "0.1.0",
        RuntimeApiVersion: "v1",
    }, nil
}


// Add minimal stub implementations for required RuntimeService methods
func (s *MyRuntime) RunPodSandbox(ctx context.Context, req *runtimeapi.RunPodSandboxRequest) (*runtimeapi.RunPodSandboxResponse, error) {
    return nil, nil
}
func (s *MyRuntime) StopPodSandbox(ctx context.Context, req *runtimeapi.StopPodSandboxRequest) (*runtimeapi.StopPodSandboxResponse, error) {
    return nil, nil
}
func (s *MyRuntime) RemovePodSandbox(ctx context.Context, req *runtimeapi.RemovePodSandboxRequest) (*runtimeapi.RemovePodSandboxResponse, error) {
    return nil, nil
}
func (s *MyRuntime) PodSandboxStatus(ctx context.Context, req *runtimeapi.PodSandboxStatusRequest) (*runtimeapi.PodSandboxStatusResponse, error) {
    return nil, nil
}
func (s *MyRuntime) ListPodSandbox(ctx context.Context, req *runtimeapi.ListPodSandboxRequest) (*runtimeapi.ListPodSandboxResponse, error) {
    return nil, nil
}
func (s *MyRuntime) CreateContainer(ctx context.Context, req *runtimeapi.CreateContainerRequest) (*runtimeapi.CreateContainerResponse, error) {
	pid, err := runContainer([]string{"tail", "/dev/null"})
	if err != nil {
		return nil, err
	}

	return &runtimeapi.CreateContainerResponse{ContainerId: string(pid)}, nil
}
func (s *MyRuntime) StartContainer(ctx context.Context, req *runtimeapi.StartContainerRequest) (*runtimeapi.StartContainerResponse, error) {
    return nil, nil
}
func (s *MyRuntime) StopContainer(ctx context.Context, req *runtimeapi.StopContainerRequest) (*runtimeapi.StopContainerResponse, error) {
    return nil, nil
}
func (s *MyRuntime) RemoveContainer(ctx context.Context, req *runtimeapi.RemoveContainerRequest) (*runtimeapi.RemoveContainerResponse, error) {
    return nil, nil
}
func (s *MyRuntime) ContainerStatus(ctx context.Context, req *runtimeapi.ContainerStatusRequest) (*runtimeapi.ContainerStatusResponse, error) {
    return nil, nil
}
func (s *MyRuntime) ListContainers(ctx context.Context, req *runtimeapi.ListContainersRequest) (*runtimeapi.ListContainersResponse, error) {
    return nil, nil
}

// --- ImageService Implementation ---

func (s *MyRuntime) ListImages(ctx context.Context, req *runtimeapi.ListImagesRequest) (*runtimeapi.ListImagesResponse, error) {
    return &runtimeapi.ListImagesResponse{}, nil
}
func (s *MyRuntime) ImageStatus(ctx context.Context, req *runtimeapi.ImageStatusRequest) (*runtimeapi.ImageStatusResponse, error) {
    return &runtimeapi.ImageStatusResponse{}, nil
}
func (s *MyRuntime) PullImage(ctx context.Context, req *runtimeapi.PullImageRequest) (*runtimeapi.PullImageResponse, error) {
    return &runtimeapi.PullImageResponse{
        ImageRef: "dummy-image-ref",
    }, nil
}
func (s *MyRuntime) RemoveImage(ctx context.Context, req *runtimeapi.RemoveImageRequest) (*runtimeapi.RemoveImageResponse, error) {
    return &runtimeapi.RemoveImageResponse{}, nil
}
func (s *MyRuntime) ImageFsInfo(ctx context.Context, req *runtimeapi.ImageFsInfoRequest) (*runtimeapi.ImageFsInfoResponse, error) {
    return &runtimeapi.ImageFsInfoResponse{}, nil
}

// Register both services to gRPC
func RegisterMyRuntime(grpcServer *grpc.Server) {
    s := &MyRuntime{}
    runtimeapi.RegisterRuntimeServiceServer(grpcServer, s)
    runtimeapi.RegisterImageServiceServer(grpcServer, s)
}
