package orchestrator

import (
	"context"

	pb "github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/orchestrator/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	pb.UnimplementedOrchestratorServer
}

func NewGRPCServer() pb.OrchestratorServer {
	return &grpcServer{}
}

func (s *grpcServer) GetTask(ctx context.Context, _ *emptypb.Empty) (*pb.TaskMessage, error) {
	return &pb.TaskMessage{Id: 1, Arg1: 2, Arg2: 3, Operation: "+"}, nil
}

func (s *grpcServer) PostResult(ctx context.Context, in *pb.TaskResult) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
