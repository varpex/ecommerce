package main

import (
	"fmt"
	"github.com/kataras/iris"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func productBaseFeaturesParty(productBaseFeatures iris.Party) {
	productBaseFeatures.Get("/", productBaseFeatureList)
	productBaseFeatures.Get("/{id: long}", productBaseFeatureRetrieve)
	productBaseFeatures.Post("/", productBaseFeaturePost)
	productBaseFeatures.Patch("/{id: long}", productBaseFeaturePatch)
	productBaseFeatures.Delete("/{id: long}", productBaseFeatureDelete)
}

func productBaseFeatureList(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instances []ProductBaseFeature
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

func productBaseFeatureRetrieve(ctx iris.Context) {
	id := ctx.Params().Get("id")
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var instance ProductBaseFeature
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	ctx.JSON(instance)
}

func productBaseFeaturePost(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	product := &ProductBaseFeature{}
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

func productBaseFeaturePatch(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	id := ctx.Params().Get("id")

	product := &ProductBaseFeature{}
	dataErr := ctx.ReadJSON(product)
	if dataErr != nil {
		panic(dataErr)
	}

	var instance ProductBaseFeature
	findErr := db.Where("id = ?", id).Find(&instance).Error
	if findErr != nil {
		panic(findErr)
	}

	instance.Title = product.Title

	saveErr := db.Save(&instance).Error
	if saveErr != nil {
		panic(saveErr)
	}

	ctx.JSON(product)
}

func productBaseFeatureDelete(ctx iris.Context) {
	db, err := gorm.Open("postgres", productConnectionString)
	panic(err)
	defer db.Close()

	id := ctx.Params().Get("id")

	var instance ProductBaseFeature
	db.Where("id = ?", id).Find(&instance)

	db.Delete(&instance)

	data := make(map[string]interface{})
	data["message"] = fmt.Sprintf("%d Deleted.", id)
	ctx.JSON(data)
}