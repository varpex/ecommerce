package main

import (
	"fmt"
	"github.com/kataras/iris"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func productFeaturesParty(productFeatures iris.Party) {
	productFeatures.Get("/", productFeatureList)
	productFeatures.Get("/{id: long}", productFeatureRetrieve)
	productFeatures.Post("/", productFeaturePost)
	productFeatures.Patch("/{id: long}", productFeaturePatch)
	productFeatures.Delete("/{id: long}", productFeatureDelete)
}

func productFeatureList(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instances []ProductFeature
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

func productFeatureRetrieve(ctx iris.Context) {
	id := ctx.Params().Get("id")
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance ProductFeature
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	ctx.JSON(instance)
}

func productFeaturePost(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	product := &ProductFeature{}
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

func productFeaturePatch(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	id := ctx.Params().Get("id")

	product := &ProductFeature{}
	dataErr := ctx.ReadJSON(product)
	if dataErr != nil {
		panic(dataErr)
	}

	var instance ProductFeature
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	instance.ParentRefer = product.ParentRefer
	instance.ProductId = product.ProductId
	instance.Value = product.Value

	saveErr := db.Save(&instance).Error
	if saveErr != nil {
		panic(saveErr)
	}

	ctx.JSON(product)
}

func productFeatureDelete(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	panic(err)
	defer db.Close()

	id := ctx.Params().Get("id")

	var instance ProductFeature
	db.Where("id = ?", id).Find(&instance)

	db.Delete(&instance)

	data := make(map[string]interface{})
	data["message"] = fmt.Sprintf("%d Deleted.", id)
	ctx.JSON(data)
}