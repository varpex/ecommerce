package main

import (
	"fmt"
	"github.com/kataras/iris"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func productPricesParty(productPrices iris.Party) {
	productPrices.Get("/", productPriceList)
	productPrices.Get("/{id: long}", productPriceRetrieve)
	productPrices.Post("/", productPricePost)
	productPrices.Patch("/{id: long}", productPricePatch)
	productPrices.Delete("/{id: long}", productPriceDelete)
}

func productPriceList(ctx iris.Context) {
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

	data := make(map[string]interface{})
	data["results"] = instances
	data["count"] = total
	// return c.RenderJSON(data)
	ctx.JSON(data)
}

func productPriceRetrieve(ctx iris.Context) {
	id := ctx.Params().Get("id")
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

	ctx.JSON(instance)
}

func productPricePost(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	product := &ProductPrice{}
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

func productPricePatch(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	id := ctx.Params().Get("id")

	product := &ProductPrice{}
	dataErr := ctx.ReadJSON(product)
	if dataErr != nil {
		panic(dataErr)
	}

	var instance ProductPrice
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	instance.Value = product.Value
	instance.ProductId = product.ProductId

	saveErr := db.Save(&instance).Error
	if saveErr != nil {
		panic(saveErr)
	}

	ctx.JSON(product)
}

func productPriceDelete(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	panic(err)
	defer db.Close()

	id := ctx.Params().Get("id")

	var instance ProductPrice
	db.Where("id = ?", id).Find(&instance)

	db.Delete(&instance)

	data := make(map[string]interface{})
	data["message"] = fmt.Sprintf("%d Deleted.", id)
	ctx.JSON(data)
}