package main

import (
	"fmt"
	"github.com/kataras/iris/mvc"

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

func productMVC(app *mvc.Application) {
    // app.Register(...)
    // app.Router.Use/UseGlobal/Done(...)
    app.Handle(new(ProductController))
}

type ProductController struct {}

func (m *ProductController) BeforeActivation(b mvc.BeforeActivation) {
    // b.Dependencies().Add/Remove
    // b.Router().Use/UseGlobal/Done // and any standard API call you already know
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		// ctx.JSON(iris.Map{"error": err.Error})
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&Product{}, &ProductPrice{}, &ProductBaseFeature{}, &ProductFeature{})
    // 1-> Method
    // 2-> Path
    // 3-> The controller's function name to be parsed as handler
    // 4-> Any handlers that should run before the MyCustomHandler
	b.Handle("GET", "/something/{id:long}", "MyCustomHandler")
	b.Handle("POST", "/", "Post")
}

// GET: http://localhost:8080/root
// func (m *ProductController) Get() string { return "Hey" }

// GET: http://localhost:8080/root/something/{id:long}
func (m *ProductController) MyCustomHandler(id int64) string { return "MyCustomHandler says Hey" }

func (m *ProductController) Get() interface{} {
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
	return data
}

func (m *ProductController) GetBy(id int64) interface{} {
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

	return instance
}

func (m *ProductController) Post(product Product) interface{} {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		fmt.Printf("%s\n", err.Error)
		// panic(err)
	}
	defer db.Close()

	// instance := &Product{
	// 	Title:       Title,
	// 	Biography:   Biography,
	// 	Description: Description,
	// }

	createErr := db.Create(product).Scan(&product).Error
	if createErr != nil {
		fmt.Printf("%s\n", createErr.Error)
		// panic(createErr)
	}

	return product
}

func (m *ProductController) PostBy(product Product) interface{} {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance Product
	findErr := db.Where("id = ?", product.ID).Find(&instance).Error
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

	saveErr := db.Save(&instance).Error
	if saveErr != nil {
		panic(saveErr)
	}

	return instance
}

func (m *ProductController) Delete(id int64) interface{} {
	db, err := gorm.Open("postgres", productConnectionString)
	panic(err)
	defer db.Close()

	var instance Product
	db.Where("id = ?", id).Find(&instance)

	db.Delete(&instance)

	data := make(map[string]interface{})
	data["message"] = fmt.Sprintf("%d Deleted.", id)
	return data
}