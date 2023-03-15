package usecase

import "github.com/raphaelmb/go-clean-architecture/internal/entity"

type OrderListOutputDTO struct {
	Orders []entity.Order
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(orderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (l *ListOrdersUseCase) Execute() (OrderListOutputDTO, error) {
	orders, err := l.OrderRepository.List()
	if err != nil {
		return OrderListOutputDTO{}, err
	}
	dto := OrderListOutputDTO{Orders: orders}

	return dto, nil
}
