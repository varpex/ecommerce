package main

import (
	"fmt"
	"strconv"

	"github.com/kataras/iris/hero"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func productPricesList() hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instances []ProductPrice
	var total int
	findErr := db.Find(&instances).Count(&total).Error
	if findErr != nil {
		panic(findErr)
	}

	var products []Product
	productsFindErr := db.Find(&products).Error
	if productsFindErr != nil {
		panic(productsFindErr)
	}

	data := make(map[string]interface{})
	data["results"] = instances
	data["count"] = total
	data["products"] = products
	data["title"] = "لیست قیمت برای محصولات"

	return hero.View{
		Name:   "productPrices/list.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productPricesListByProduct(product int64) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instances []ProductPrice
	var total int
	findErr := db.Where("product_id = ?", product).Find(&instances).Count(&total).Error
	if findErr != nil {
		panic(findErr)
	}

	var products []Product
	productsFindErr := db.Where("id = ?", product).Find(&products).Error
	if productsFindErr != nil {
		panic(productsFindErr)
	}

	data := make(map[string]interface{})
	data["results"] = instances
	data["count"] = total
	data["products"] = products
	data["title"] = "تصاویر برای محصول: " + strconv.FormatInt(product, 10)

	return hero.View{
		Name:   "productPrices/list.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productPricesRetrieve(id int64) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance ProductPrice
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	var products []Product
	productsFindErr := db.Find(&products).Error
	if productsFindErr != nil {
		panic(productsFindErr)
	}

	data := make(map[string]interface{})
	data["instance"] = instance
	data["title"] = "Product Price"
	data["products"] = products

	return hero.View{
		Name:   "productPrices/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productPricesPost(productPrice ProductPrice) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createErr := db.Create(&productPrice).Scan(&productPrice).Error
	if createErr != nil {
		panic(createErr)
	}

	var products []Product
	productsFindErr := db.Find(&products).Error
	if productsFindErr != nil {
		panic(productsFindErr)
	}

	data := make(map[string]interface{})
	data["instance"] = productPrice
	data["title"] = "Product Price Detail"
	data["products"] = products

	return hero.View{
		Name:   "productPrices/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productPricesPatch(id int64, productPrice ProductPrice) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance ProductPrice
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	if productPrice.Value != 0 {
		instance.Value = productPrice.Value
	}
	if productPrice.ProductId != 0 {
		instance.ProductId = productPrice.ProductId
	}

	saveErr := db.Save(&instance).Error
	if saveErr != nil {
		panic(saveErr)
	}

	var products []Product
	productsFindErr := db.Find(&products).Error
	if productsFindErr != nil {
		panic(productsFindErr)
	}

	data := make(map[string]interface{})
	data["instance"] = instance
	data["title"] = "ProductPrice Detail"
	data["products"] = products

	return hero.View{
		Name:   "productPrices/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productPricesDelete(id int64) interface{} {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance ProductPrice
	db.Where("id = ?", id).Find(&instance)

	db.Delete(&instance)

	data := make(map[string]interface{})
	data["message"] = fmt.Sprintf("%d Deleted.", id)
	return data
}
