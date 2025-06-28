package usecase

import (
	"clean-architecture/internal/dto"
	"clean-architecture/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrderUseCase) Execute() ([]dto.OrderOutputDTO, error) {
	var orders = []dto.OrderOutputDTO{}

	ordersEntity, err := c.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	for _, order := range ordersEntity {
		dto := dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
		orders = append(orders, dto)
	}

	return orders, nil
}
