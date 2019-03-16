package main

import (
	"fmt"
	"strconv"

	"github.com/kataras/iris/hero"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func productImagesList() hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instances []ProductImage
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
	data["title"] = "تصاویر برای محصولات"

	return hero.View{
		Name:   "productImages/list.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productImagesListByProduct(product int64) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instances []ProductImage
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
		Name:   "productImages/list.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productImagesRetrieve(id int64) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance ProductImage
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
	data["title"] = "Product Image"
	data["products"] = products

	return hero.View{
		Name:   "productImages/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productImagesPost(productImage ProductImage) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createErr := db.Create(&productImage).Scan(&productImage).Error
	if createErr != nil {
		panic(createErr)
	}

	var products []Product
	productsFindErr := db.Find(&products).Error
	if productsFindErr != nil {
		panic(productsFindErr)
	}

	data := make(map[string]interface{})
	data["instance"] = productImage
	data["title"] = "ProductImage Detail"
	data["products"] = products

	return hero.View{
		Name:   "productImages/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productImagesPatch(id int64, productImage ProductImage) hero.Result {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance ProductImage
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	if productImage.Path != "" {
		instance.Path = productImage.Path
	}
	if productImage.ProductId != 0 {
		instance.ProductId = productImage.ProductId
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
	data["title"] = "ProductImage Detail"
	data["products"] = products

	return hero.View{
		Name:   "productImages/form.html",
		Layout: "admin/main.html",
		Data:   data,
	}
}

func productImagesDelete(id int64) interface{} {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance ProductImage
	db.Where("id = ?", id).Find(&instance)

	db.Delete(&instance)

	data := make(map[string]interface{})
	data["message"] = fmt.Sprintf("%d Deleted.", id)
	return data
}
