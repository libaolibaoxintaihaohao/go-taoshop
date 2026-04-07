package repository

import (
	"context"
	"database/sql"

	"taoshop/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {//查询商品列表
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, price, stock, cover_url, description, created_at
		FROM products
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {//循环下一行数据，如果有下一行数据就继续循环
		var product models.Product
		if err := rows.Scan(
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
		products = append(products, product)//把刚才填好的这个 product 对象，加到 products 这个大列表的末尾。
	}

	return products, rows.Err()//检查一下循环过程中有没有发生错误
}

func (r *ProductRepository) FindByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*models.Product, error) {//事务中按ID查找
	row := tx.QueryRowContext(ctx, `
		SELECT id, name, price, stock, cover_url, description, created_at
		FROM products
		WHERE id = ?
		FOR UPDATE`, id)

	var product models.Product
	if err := row.Scan(
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

	return &product, nil
}

func (r *ProductRepository) ReduceStockTx(ctx context.Context, tx *sql.Tx, productID int64, quantity int) error {//事务中扣减库存
	_, err := tx.ExecContext(ctx, `
		UPDATE products
		SET stock = stock - ?
		WHERE id = ?`, quantity, productID)
	return err
}
