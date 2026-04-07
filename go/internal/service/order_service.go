package service

import (
	"context"
	"database/sql"
	"errors"

	"taoshop/internal/models"
	"taoshop/internal/repository"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type OrderService struct {
	db       *sql.DB
	products *repository.ProductRepository
	orders   *repository.OrderRepository
	catalog  *CatalogService
}

func NewOrderService(db *sql.DB, products *repository.ProductRepository, orders *repository.OrderRepository, catalog *CatalogService) *OrderService {
	return &OrderService{
		db:       db,
		products: products,
		orders:   orders,
		catalog:  catalog,
	}
}

func (s *OrderService) Create(ctx context.Context, userID, productID int64, quantity int) (*models.Order, error) {
	// 下单是典型的事务场景：查库存、扣库存、写订单，必须同时成功或同时失败。
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	product, err := s.products.FindByIDTx(ctx, tx, productID)
	if err != nil {
		return nil, err
	}

	if product.Stock < quantity {
		return nil, ErrInsufficientStock
	}

	order := &models.Order{
		UserID:     userID,
		ProductID:  productID,
		Product:    *product,
		Quantity:   quantity,
		TotalPrice: product.Price * float64(quantity),
		Status:     "paid",
	}

	if err := s.products.ReduceStockTx(ctx, tx, productID, quantity); err != nil {
		return nil, err
	}

	if err := s.orders.CreateTx(ctx, tx, order); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 商品数据变了以后，主动删缓存，下一次查询会重新回源 MySQL。
	s.catalog.InvalidateProducts(ctx)
	return order, nil
}

func (s *OrderService) ListByUser(ctx context.Context, userID int64) ([]models.Order, error) {
	return s.orders.ListByUser(ctx, userID)
}
