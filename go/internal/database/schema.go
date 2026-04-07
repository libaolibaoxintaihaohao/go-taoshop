package database

import "database/sql"

func EnsureSchema(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(64) NOT NULL UNIQUE,
			email VARCHAR(128) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS products (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(120) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			stock INT NOT NULL,
			cover_url VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS orders (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			product_id BIGINT NOT NULL,
			quantity INT NOT NULL,
			total_price DECIMAL(10,2) NOT NULL,
			status VARCHAR(32) NOT NULL DEFAULT 'paid',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id),
			CONSTRAINT fk_orders_product FOREIGN KEY (product_id) REFERENCES products(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`INSERT INTO products (name, price, stock, cover_url, description)
		SELECT * FROM (
			SELECT '轻薄笔记本 Pro', 6999.00, 48, 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?auto=format&fit=crop&w=800&q=80', '适合办公与轻度创作的高性能轻薄本'
			UNION ALL
			SELECT '无线降噪耳机 Max', 1299.00, 120, 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?auto=format&fit=crop&w=800&q=80', '长续航、低延迟，适合通勤与影音娱乐'
			UNION ALL
			SELECT '人体工学办公椅', 1899.00, 36, 'https://images.unsplash.com/photo-1505843513577-22bb7d21e455?auto=format&fit=crop&w=800&q=80', '支撑腰背，适合久坐学习与工作'
			UNION ALL
			SELECT '电竞显示器 27 寸', 2399.00, 25, 'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?auto=format&fit=crop&w=800&q=80', '2K 高刷，适合游戏与多任务场景'
		) AS seed
		WHERE NOT EXISTS (SELECT 1 FROM products LIMIT 1);`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {//执行stmt里的sql语句，错了就返回
			return err
		}
	}

	return nil
}
