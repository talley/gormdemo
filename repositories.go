package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name      string
	Price     float64
	Available bool
}

func GetProducts() ([]Product, error) {
	db, err, ctx := InitializeDb()
	if err != nil {
		fmt.Println(err)
	}
	var products []Product
	if err := db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, err
}

func AddProduct(product Product) error {
	db, err, ctx := InitializeDb()
	err = gorm.G[Product](db).Create(ctx, &Product{Name: product.Name, Price: product.Price, Available: product.Available})
	if err != nil {
		fmt.Println("Failed to create product")
		fmt.Println(err)
	}
	return err
}
func DeleteProduct0(productId int) (int, error) {
	db, err, ctx := InitializeDb()
	if err != nil {
		fmt.Println(err)
	}
	i, err := gorm.G[Product](db).Where("id = ?", productId).Delete(ctx)
	if err != nil {
		fmt.Println("Failed to delete product")
		fmt.Println(err)
	}
	return i, err
}
func DeleteProduct3(productId int) error {
	db, err, ctx := InitializeDb()
	ctxerror := ctx.Err()
	if ctxerror != nil {
		fmt.Println(ctxerror)
	}
	if err != nil {

		fmt.Errorf("failed to initialize database: %w", err)
	}
	db.Delete(&Product{}, productId)
	// Ensure database connection is closed
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	fmt.Println("Successfully deleted product")
	return err
}
func DeleteProduct(productId int) (int, error) {
	db, err, ctx := InitializeDb()
	if err != nil {
		return 0, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Ensure database connection is closed
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Method 1: Using Delete with model instance
	//err := gorm.G[Email](db).Where("id = ?", 10).Delete(ctx)
	result := db.WithContext(ctx).Where("id = ?", productId).Delete(&Product{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete product: %w", result.Error)
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("product with id %d not found", productId)
	}

	return int(result.RowsAffected), nil
}

func DeleteProductAlternative(productId int) (int, error) {
	db, err, ctx := InitializeDb()
	if err != nil {
		return 0, fmt.Errorf("failed to initialize database: %w", err)
	}

	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// First find the product to ensure it exists
	var product Product
	if err := db.WithContext(ctx).First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("product with id %d not found", productId)
		}
		return 0, fmt.Errorf("failed to find product: %w", err)
	}

	// Then delete it
	result := db.WithContext(ctx).Delete(&product)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete product: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}
func UpdateProduct(id int, product Product) (Product, error) {
	db, err, ctx := InitializeDb()
	if err != nil {
		fmt.Println(err)
	}
	// Ensure the ID is set
	if id == 0 {
		return Product{}, fmt.Errorf("missing product ID")
	}
	findProduct, err := GetProduct(id)
	if err != nil {
		fmt.Println(err)
	}

	// Use a map so zero-values (e.g. false, 0, "") are written
	res := db.WithContext(ctx).
		Model(&Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]interface{}{
			"name":       product.Name,
			"price":      product.Price,
			"available":  product.Available,
			"updated_at": time.Now(),
			"created_at": findProduct.CreatedAt,
			"deleted_at": findProduct.DeletedAt,
		})
	if res.Error != nil {
		return Product{}, res.Error
	}
	if res.RowsAffected == 0 {
		return Product{}, gorm.ErrRecordNotFound
	}

	// Optionally fetch the updated row
	var updated Product
	if err := db.WithContext(ctx).First(&updated, product.ID).Error; err != nil {
		return Product{}, err
	}
	return updated, nil
}
func GetProduct(id int) (Product, error) {
	db, err, ctx := InitializeDb()
	product, err := gorm.G[Product](db).Where("id = ?", id).First(ctx)
	if err != nil {
		fmt.Println("Failed to get product")
		fmt.Println(err)
	}
	return product, err
}
func InitializeDb() (*gorm.DB, error, context.Context) {
	dsn := "host=localhost user=postgres password=Iamsmart27! dbname=gormdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	ctx := context.Background()

	db.AutoMigrate(&Product{})
	return db, err, ctx
}
