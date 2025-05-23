package entity

import "time"

type Product struct {
	ID                  int       `db:"id"`
	Name                string    `db:"name"`
	RealPrice           int       `db:"real_price"`
	Description         string    `db:"description"`
	Stock               int       `db:"stock"`
	ImageUrl            string    `db:"image_url"`
	IsActive            bool      `db:"is_active"`
	ProductCategoryID   int       `db:"product_category_id"`
	ProductCategoryName string    `db:"product_category_name"`
	RecipeID            int       `db:"recipe_id"`
	CreatedAt           time.Time `db:"created_at"`
}
