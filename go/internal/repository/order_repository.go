package repository

import (
	"context"
	"database/sql"

	"taoshop/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateTx(ctx context.Context, tx *sql.Tx, order *models.Order) error {//创建订单
	result, err := tx.ExecContext(ctx, `
		INSERT INTO orders (user_id, product_id, quantity, total_price, status)
		VALUES (?, ?, ?, ?, ?)`,
		order.UserID, order.ProductID, order.Quantity, order.TotalPrice, order.Status,
	)
	if err != nil {
		return err
	}

	order.ID, _ = result.LastInsertId()//自动生成一个唯一的订单号（ID）,然后行代码把这个id取出来
	return nil
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID int64) ([]models.Order, error) {//查出某个人的所有订单，并且把订单里对应的商品信息也一起查出来。
	rows, err := r.db.QueryContext(ctx, `//把符合条件的行全部拿来
		SELECT 
			o.id, o.user_id, o.product_id, o.quantity, o.total_price, o.status, o.created_at,
			p.id, p.name, p.price, p.stock, p.cover_url, p.description, p.created_at
		FROM orders o
		INNER JOIN products p ON p.id = o.product_id
		WHERE o.user_id = ?
		ORDER BY o.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var product models.Product

		if err := rows.Scan(//把查出来的数据都填进去
			&order.ID,
			&order.UserID,
			&order.ProductID,
			&order.Quantity,
			&order.TotalPrice,
			&order.Status,
			&order.CreatedAt,
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.CoverURL,
			&product.Description,
			&product.CreatedAt,
		); err != nil {
			return nil, err
		}

		order.Product = product//包装一下
		orders = append(orders, order)
	}

	return orders, rows.Err()
}
