package main

import (
	"fmt"

	"github.com/kataras/iris/hero"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ProductCategory struct {
	gorm.Model
	Title string `gorm:"index:idx_product_categories_title"`
}

type Product struct {
	gorm.Model
	Title         string           `gorm:"index:idx_products_title"`
	Biography     string           `gorm:"index:idx_products_biography"`
	Description   string           `gorm:"index:idx_products_description"`
	Url           string           `gorm:"index:idx_products_url"`
	Category      ProductCategory  `gorm:"foreignkey:CategoryRefer"`
	CategoryRefer uint             `gorm:"index:idx_products_category_refer"`
	Prices        []ProductPrice   `gorm:"ForeignKey:ProductId"`
	Features      []ProductFeature `gorm:"ForeignKey:ProductId"`
}

type ProductPrice struct {
	gorm.Model
	Value     uint `gorm:"index:idx_product_prices_value"`
	ProductId uint `gorm:"index:idx_product_prices_product_id"`
}

type ProductImage struct {
	gorm.Model
	Path      string `gorm:"index:idx_product_images_path"`
	ProductId uint   `gorm:"index:idx_product_images_product_id"`
}

type ProductBaseFeature struct {
	gorm.Model
	Title string `gorm:"index:idx_product_base_features_title"`
}

type ProductFeature struct {
	gorm.Model
	Feature     ProductBaseFeature `gorm:"foreignkey:ParentRefer"`
	ParentRefer uint               `gorm:"index:idx_product_features_parent_refer"`
	ProductId   uint               `gorm:"index:idx_product_features_product_id"`
	Value       string             `gorm:"index:idx_product_features_value"`
}

var (
	productConnectionString = "host=185.105.239.12 user=postgres password=1369s1r3d691369 dbname=products sslmode=disable"
)

func productsList() hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instances []Product
	var total int
	var categories []ProductCategory

	findErr := db.Preload("Prices").Find(&instances).Count(&total).Error
	if findErr != nil {
		panic(findErr)
	}

	categoryFindErr := db.Find(&categories).Error
	if categoryFindErr != nil {
		panic(categoryFindErr)
	}

	data := make(map[string]interface{})
	data["results"] = instances
	data["count"] = total
	data["categories"] = categories

	return hero.View{
		Name:   "products/list.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productsRetrieve(id int64) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance Product
	findErr := db.Where("id = ?", id).Preload("Prices").Preload("Features").Preload("Category").Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	var categories []ProductCategory
	categoryFindErr := db.Find(&categories).Error
	if categoryFindErr != nil {
		panic(categoryFindErr)
	}

	data := make(map[string]interface{})
	data["instance"] = instance
	data["title"] = instance.Title
	data["categories"] = categories

	return hero.View{
		Name:   "products/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productsPost(product Product) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createErr := db.Create(&product).Scan(&product).Error
	if createErr != nil {
		panic(createErr)
	}

	data := make(map[string]interface{})
	data["instance"] = product
	data["title"] = "Product Detail"

	return hero.View{
		Name:   "products/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func uploadList() interface{} {
	return hero.View{
		Name:   "upload.html",
		Layout: "admin/main.html",
	}
}

func uploadPost(upload ProductImage) interface{} {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return upload
}

func productsPatch(id int64, product Product) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance Product
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	if product.Title != "" {
		instance.Title = product.Title
	}
	if product.Biography != "" {
		instance.Biography = product.Biography
	}
	if product.Description != "" {
		instance.Description = product.Description
	}
	if product.Url != "" {
		instance.Url = product.Url
	}
	if product.CategoryRefer != 0 {
		instance.CategoryRefer = product.CategoryRefer
	}

	saveErr := db.Save(&instance).Error
	if saveErr != nil {
		panic(saveErr)
	}

	data := make(map[string]interface{})
	data["instance"] = instance
	data["title"] = "Product Detail"

	return hero.View{
		Name:   "products/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productsDelete(id int64) interface{} {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance Product
	db.Where("id = ?", id).Find(&instance)

	db.Delete(&instance)

	data := make(map[string]interface{})
	data["message"] = fmt.Sprintf("%d Deleted.", id)
	return data
}
