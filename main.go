// package main

// import (
// 	"github.com/kataras/iris"

// 	"github.com/kataras/iris/middleware/logger"
// 	"github.com/kataras/iris/middleware/recover"
// )

// func main() {
// 	app := iris.New()
// 	app.Logger().SetLevel("debug")
// 	// Optionally, add two built'n handlers
// 	// that can recover from any http-relative panics
// 	// and log the requests to the terminal.
// 	app.Use(recover.New())
// 	app.Use(logger.New())

// 	app.RegisterView(iris.HTML("./templates", ".html"))

// 	// Method:   GET
// 	// Resource: http://localhost:8080
// 	app.Handle("GET", "/", func(ctx iris.Context) {
// 		ctx.ViewLayout("main.html")
// 		ctx.ViewData("title", "Home")
// 		ctx.ViewData("name", "Saeid Ramezani")
// 		ctx.View("index.html")
// 	})

// 	// same as app.Handle("GET", "/ping", [...])
// 	// Method:   GET
// 	// Resource: http://localhost:8080/ping
// 	app.Post("/ping", func(ctx iris.Context) {
// 		ctx.WriteString(ctx.FormValue("name"))
// 	})

// 	// Method:   GET
// 	// Resource: http://localhost:8080/hello
// 	app.Get("/hello", func(ctx iris.Context) {
// 		ctx.JSON(iris.Map{"message": "Hello Iris!"})
//     })

// 		app.PartyFunc("/products", productsParty)
// 		app.PartyFunc("/products/prices/", productPricesParty)
// 		app.PartyFunc("/products/base-features/", productBaseFeaturesParty)
// 		app.PartyFunc("/products/features/", productFeaturesParty)

// 	// http://localhost:8080
// 	// http://localhost:8080/ping
// 	// http://localhost:8080/hello
// 	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
// }

package main

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	app := iris.New()

	db, err := gorm.Open("postgres", productConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&ProductCategory{}, &Product{}, &ProductPrice{}, &ProductImage{}, &ProductBaseFeature{}, &ProductFeature{})

	app.RegisterView(iris.HTML("./templates", ".html"))

	hero.Register(func(ctx iris.Context) (product Product) {
		ctx.ReadForm(&product)
		return
	})

	hero.Register(func(ctx iris.Context) (productPrice ProductPrice) {
		ctx.ReadForm(&productPrice)
		return
	})

	hero.Register(func(ctx iris.Context) (productBaseFeature ProductBaseFeature) {
		ctx.ReadForm(&productBaseFeature)
		return
	})

	hero.Register(func(ctx iris.Context) (productFeature ProductFeature) {
		ctx.ReadForm(&productFeature)
		return
	})

	hero.Register(func(ctx iris.Context) (productImage ProductImage) {
		ctx.ReadForm(&productImage)

		size, err := ctx.UploadFormFiles("./assets/uploads", beforeSave)
		if err != nil {
			fmt.Println(err.Error)
		}
		if size != 0 {
			data := ctx.GetViewData()
			instance := &productImage
			instance.Path = "/uploads/" + data["fileName"].(string)
		}
		return
	})

	// Product
	productsListHandler := hero.Handler(productsList)
	app.Get("/products", productsListHandler)

	productsRetrieveHandler := hero.Handler(productsRetrieve)
	app.Get("/products/{id: long min(1)}", productsRetrieveHandler)

	productsPostHandler := hero.Handler(productsPost)
	app.Post("/products", productsPostHandler)

	productsPatchHandler := hero.Handler(productsPatch)
	app.Patch("/products/{id: long min(1)}", productsPatchHandler)
	app.Post("/products/{id: long min(1)}", productsPatchHandler)

	productsDeleteHandler := hero.Handler(productsDelete)
	app.Delete("/products/{id: long min(1)}", productsDeleteHandler)

	// Product Image
	productImagesListHandler := hero.Handler(productImagesList)
	app.Get("/product-images", productImagesListHandler)

	productImagesByProductListHandler := hero.Handler(productImagesListByProduct)
	app.Get("/products/{product: long min(1)}/images", productImagesByProductListHandler)

	productImagesRetrieveHander := hero.Handler(productImagesRetrieve)
	app.Get("/product-images/{id: long min(1)}", productImagesRetrieveHander)

	productImagesPostHandler := hero.Handler(productImagesPost)
	app.Post("/product-images", productImagesPostHandler)

	productImagesPatchHandler := hero.Handler(productImagesPatch)
	app.Patch("/product-images/{id: long min(1)}", productImagesPatchHandler)
	app.Post("/product-images/{id: long min(1)}", productImagesPatchHandler)

	productImagesDeleteHandler := hero.Handler(productImagesDelete)
	app.Delete("/product-images/{id: long min(1)}", productImagesDeleteHandler)

	// Product Price
	productPricesListHandler := hero.Handler(productPricesList)
	app.Get("/product-prices", productPricesListHandler)

	productPricesByProductListHandler := hero.Handler(productPricesListByProduct)
	app.Get("/products/{product: long min(1)}/prices", productPricesByProductListHandler)

	productPricesRetrieveHander := hero.Handler(productPricesRetrieve)
	app.Get("/product-prices/{id: long min(1)}", productPricesRetrieveHander)

	productPricesPostHandler := hero.Handler(productPricesPost)
	app.Post("/product-prices", productPricesPostHandler)

	productPricesPatchHandler := hero.Handler(productPricesPatch)
	app.Patch("/product-prices/{id: long min(1)}", productPricesPatchHandler)
	app.Post("/product-prices/{id: long min(1)}", productPricesPatchHandler)

	productPricesDeleteHandler := hero.Handler(productPricesDelete)
	app.Delete("/product-prices/{id: long min(1)}", productPricesDeleteHandler)

	app.StaticWeb("/static", "./assets")
	// 1
	helloHandler := hero.Handler(hello)
	app.Get("/{to:string}", helloHandler)

	// 2
	hero.Register(&myTestService{
		prefix: "Service: Hello",
	})

	helloServiceHandler := hero.Handler(helloService)
	app.Get("/service/{to:string}", helloServiceHandler)

	// 3
	hero.Register(func(ctx iris.Context) (form LoginForm) {
		// it binds the "form" with a
		// x-www-form-urlencoded form data and returns it.
		ctx.ReadForm(&form)
		return
	})

	loginHandler := hero.Handler(login)
	app.Post("/login", loginHandler)

	// http://localhost:8080/your_name
	// http://localhost:8080/service/your_name
	app.Run(
		iris.Addr(":8080"),
		iris.WithOptimizations,
		iris.WithPostMaxMemory(32<<20),
	)
}

func beforeSave(ctx iris.Context, file *multipart.FileHeader) {
	ip := ctx.RemoteAddr()
	// make sure you format the ip in a way
	// that can be used for a file name (simple case):
	ip = strings.Replace(ip, ".", "_", -1)
	ip = strings.Replace(ip, ":", "_", -1)

	// you can use the time.Now, to prefix or suffix the files
	// based on the current time as well, as an exercise.
	// i.e unixTime :=	time.Now().Unix()
	// prefix the Filename with the $IP-
	// no need for more actions, internal uploader will use this
	// name to save the file into the "./uploads" folder.
	unixTime := time.Now().Unix()
	file.Filename = ip + "-" + strconv.FormatInt(unixTime, 10) + "-" + file.Filename
	ctx.ViewData("fileName", file.Filename)
}

func hello(to string) string {
	return "Hello " + to
}

type Service interface {
	SayHello(to string) string
}

type myTestService struct {
	prefix string
}

func (s *myTestService) SayHello(to string) string {
	return s.prefix + " " + to
}

func helloService(to string, service Service) string {
	return service.SayHello(to)
}

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func login(form LoginForm) string {
	return "Hello " + form.Username
}
