package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"taoshop/internal/models"
	"taoshop/internal/repository"
)

type CatalogService struct {
	products *repository.ProductRepository
	redis    *redis.Client
}

func NewCatalogService(products *repository.ProductRepository, redis *redis.Client) *CatalogService {
	return &CatalogService{products: products, redis: redis}
}

func (s *CatalogService) ListProducts(ctx context.Context) ([]models.Product, error) {
	// 先查 Redis，命中则直接返回，这是最常见的缓存旁路模式。
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, "catalog:products").Result()
		if err == nil {
			var products []models.Product
			if json.Unmarshal([]byte(cached), &products) == nil {
				return products, nil
			}
		}
	}

	products, err := s.products.List(ctx)
	if err != nil {
		return nil, err
	}

	// 没命中缓存时回源数据库，并写回 Redis，减少热点查询直接打到 MySQL。
	if s.redis != nil {
		if payload, err := json.Marshal(products); err == nil {
			_ = s.redis.Set(ctx, "catalog:products", payload, 2*time.Minute).Err()
		}
	}

	return products, nil
}

func (s *CatalogService) InvalidateProducts(ctx context.Context) {
	if s.redis != nil {
		_ = s.redis.Del(ctx, "catalog:products").Err()
	}
}
