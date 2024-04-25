package order

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"order/internal/inventory"
	"order/internal/notification"
	"order/internal/payment"
	"order/internal/shipment"
)

type Service interface {
	AddItemToOrder(ctx context.Context, userId string, item *Item) (*Order, error)
	FinishOrder(ctx context.Context, request *FinishOrder) error
}

type service struct {
	repository   Repository
	inventory    inventory.Client
	payment      payment.Client
	notification notification.Client
	shipment     shipment.Client
}

func (s *service) AddItemToOrder(ctx context.Context, userId string, item *Item) (*Order, error) {

	order, err := s.repository.ExistOrderByStatus(ctx, userId, PENDING)
	if err != nil {
		return nil, err
	}

	if order == nil {
		order = &Order{
			UserID:      userId,
			Status:      PENDING,
			TotalAmount: decimal.Zero,
		}

		err = s.repository.Create(ctx, order)
		if err != nil {
			return nil, err
		}
	}

	item.UnitPrice, err = s.inventory.FindItemPrice(ctx, item.ProductID)
	if err != nil {
		return nil, err
	}

	var itemAlreadyExist bool
	for i := range order.Items {
		if order.Items[i].ProductID == item.ProductID && order.Items[i].UnitPrice.Equal(item.UnitPrice) {
			order.Items[i].Quantity += item.Quantity
			itemAlreadyExist = true
		}
	}

	if !itemAlreadyExist {
		order.Items = append(order.Items, *item)
	}

	totalItemPrice := item.UnitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))

	order.TotalAmount = order.TotalAmount.Add(totalItemPrice)

	err = s.repository.Update(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) FinishOrder(ctx context.Context, request *FinishOrder) error {

	order, err := s.repository.FindByOrderId(ctx, request.OrderID)
	if err != nil {
		return err
	}

	if order.Status != PENDING {
		return errors.New("order is not pending")
	}

	products := make([]inventory.Product, len(order.Items))
	for i := range order.Items {
		products[i] = inventory.Product{
			ID:       order.Items[i].ProductID,
			Quantity: order.Items[i].Quantity,
		}
	}

	err = s.inventory.CheckStock(ctx, products)
	if err != nil {
		return err
	}

	err = s.inventory.RemoveItemsFromStock(ctx, products)
	if err != nil {
		return err
	}

	err = s.payment.SendPayment(ctx)
	if err != nil {
		return err
	}

	order.Status = PROCESSING
	err = s.repository.Update(ctx, order)
	if err != nil {
		return err
	}

	err = s.notification.SendMessage(ctx, notification.Message{
		UserID:  order.UserID,
		OrderId: order.ID,
		Message: "Your order has been placed",
	})
	if err != nil {
		return err
	}

	sendItems := make([]shipment.Item, len(order.Items))
	for i := range order.Items {
		sendItems[i] = shipment.Item{
			ProductId: order.Items[i].ProductID,
			Quantity:  order.Items[i].Quantity,
		}
	}

	err = s.shipment.SendItems(ctx, sendItems, request.OrderID)
	if err != nil {
		return err
	}

	err = s.notification.SendMessage(ctx, notification.Message{
		UserID:  order.UserID,
		OrderId: order.ID,
		Message: "Your order has been shipped",
	})
	if err != nil {
		return err
	}

	order.Status = FINISHED
	err = s.repository.Update(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func NewService(repository Repository,
	inventory inventory.Client,
	payment payment.Client,
	notification notification.Client,
	shipment shipment.Client,
) Service {
	return &service{
		repository:   repository,
		inventory:    inventory,
		payment:      payment,
		notification: notification,
		shipment:     shipment,
	}
}