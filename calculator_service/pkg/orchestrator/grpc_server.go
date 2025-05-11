package orchestrator

import (
	"context"

	pb "github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/orchestrator/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcOrchServer struct {
	pb.UnimplementedOrchestratorServer
}

func NewGRPCServer() pb.OrchestratorServer {
	return &grpcOrchServer{}
}

func (s *grpcOrchServer) GetTask(ctx context.Context, _ *emptypb.Empty) (*pb.TaskMessage, error) {
	task, err := FetchNextPendingTask(ctx)
	if err != nil {
		if err == ErrNoTask {
			return nil, err
		}
		return nil, err
	}
	return &pb.TaskMessage{
		Id:            uint32(task.ID),
		Arg1:          task.Arg1,
		Arg2:          task.Arg2,
		Operation:     task.Operation,
		OperationTime: uint32(task.OperationTime),
	}, nil
}

func (s *grpcOrchServer) PostResult(ctx context.Context, in *pb.TaskResult) (*emptypb.Empty, error) {
	if err := StoreTaskResult(ctx, uint(in.Id), in.Result); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
