package service

import (
	"context"

	"clean-architecture/internal/dto"
	"clean-architecture/internal/infra/grpc/pb"
	"clean-architecture/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	dto := dto.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
		CreatedAt:  output.CreatedAt,
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *emptypb.Empty) (*pb.ListOrdersResponse, error) {
	output, err := s.ListOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	orders := make([]*pb.OrderResponse, len(output))
	for i, o := range output {
		orders[i] = &pb.OrderResponse{
			Id:         o.ID,
			Price:      float32(o.Price),
			Tax:        float32(o.Tax),
			FinalPrice: float32(o.FinalPrice),
			CreatedAt:  o.CreatedAt,
		}
	}

	if len(orders) == 0 {
		orders = []*pb.OrderResponse{}
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}
