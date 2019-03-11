package main

import (
	"fmt"
	"github.com/kataras/iris"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Product struct {
	gorm.Model
	Title       string      `gorm:"index:idx_products_title"`
	Biography   string      `gorm:"index:idx_products_biography"`
	Description string      `gorm:"index:idx_products_description"`
	Prices      []ProductPrice   `gorm:"ForeignKey:ProductId"`
	Features    []ProductFeature `gorm:"ForeignKey:ProductId"`
}

type ProductPrice struct {
	gorm.Model
	Value uint `gorm:"index:idx_product_prices_value"`
	ProductId  uint `gorm:"index:idx_product_prices_product_id"`
}

type ProductBaseFeature struct {
	gorm.Model
	Title string `gorm:"index:idx_product_base_features_title"`
}

type ProductFeature struct {
	gorm.Model
	Feature     ProductBaseFeature `gorm:"foreignkey:ParentRefer"`
	ParentRefer uint          `gorm:"index:idx_product_features_parent_refer"`
	ProductId        uint          `gorm:"index:idx_product_features_product_id"`
	Value       string        `gorm:"index:idx_product_features_value"`
}

var (
	productConnectionString = "host=185.105.239.12 user=postgres password=1369s1r3d691369 dbname=products sslmode=disable"
)

func productsParty(products iris.Party) {
	products.Get("/", List)
	products.Get("/{id: long}", Retrieve)
	products.Post("/", Post)
	products.Patch("/{id: long}", Patch)
	products.Delete("/{id: long}", Delete)
}

func List(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instances []Product
	var total int

	findErr := db.Preload("Prices").Find(&instances).Count(&total).Error
	if findErr != nil {
		panic(findErr)
	}

	data := make(map[string]interface{})
	data["results"] = instances
	data["count"] = total
	// return c.RenderJSON(data)
	ctx.JSON(data)
}

func Retrieve(ctx iris.Context) {
	id := ctx.Params().Get("id")
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance Product
	findErr := db.Where("id = ?", id).Preload("Prices").Preload("Features").Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	ctx.JSON(instance)
}

func Post(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	product := &Product{}
	dataErr := ctx.ReadJSON(product)
	if dataErr != nil {
		panic(dataErr)
	}

	createErr := db.Create(product).Scan(&product).Error
	if createErr != nil {
		panic(createErr)
	}

	ctx.JSON(product)
}

func Patch(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	id := ctx.Params().Get("id")

	product := &Product{}
	dataErr := ctx.ReadJSON(product)
	if dataErr != nil {
		panic(dataErr)
	}

	var instance Product
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	instance.Title = product.Title
	instance.Biography = product.Biography
	instance.Description = product.Description

	saveErr := db.Save(&instance).Error
	if saveErr != nil {
		panic(saveErr)
	}

	ctx.JSON(product)
}

func Delete(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	panic(err)
	defer db.Close()

	id := ctx.Params().Get("id")

	var instance Product
	db.Where("id = ?", id).Find(&instance)

	db.Delete(&instance)

	data := make(map[string]interface{})
	data["message"] = fmt.Sprintf("%d Deleted.", id)
	ctx.JSON(data)
}