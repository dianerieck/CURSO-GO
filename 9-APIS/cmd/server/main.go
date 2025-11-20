package main

import (
	"net/http"

	"github.com/dianerieck/CURSO-GO/9-APIS/configs"
	_ "github.com/dianerieck/CURSO-GO/9-APIS/docs"
	"github.com/dianerieck/CURSO-GO/9-APIS/internal/entity"
	"github.com/dianerieck/CURSO-GO/9-APIS/internal/infra/database"
	"github.com/dianerieck/CURSO-GO/9-APIS/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go Expert Example API
// @version         1.0
// @description     Product API with authentication.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Diane Elo Rieck
// @contact.url    http://www.dianerieck.com.br
// @contact.email  dianerieck@hotmail.com

// @license.name  Diane Rieck Licence
// @license.url   http://www.dianerieck.com.br

// @host      localhost:8080
// @BasePath  /

// @SecurityDefinitions.ApiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JWTExperesIn", configs.JWTExperesIn))
	//r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)

	})
	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))
	http.ListenAndServe(":8080", r)

}

//func LogRequest(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		r.Context().Value("user")
//		log.Println(r.Method, r.URL.Path)
//		next.ServeHTTP(w, r)
//	})
//}

//usar swag init no terminal para iniciar o swager e criar a pasta docs
//http://localhost:8080/docs/index.html
