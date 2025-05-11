package main

// go build && project_1.exe
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"project_1/internal/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger" // HTTP middleware for Swagger UI

	// Swagger embed files
	_ "project_1/docs" // Import the generated Swagger docs
)

type apiConfig struct {
	DB *database.Queries
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
func main() {

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}
	fmt.Println("Server is running on port: ", port)

	db_url := os.Getenv("DB_URL")
	log.Printf("DB_URL: %v", db_url)
	if db_url == "" {
		log.Fatal("DB_URL is not set")
	}

	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	db := database.New(conn)
	// Create a new instance of the API
	apiCfg := apiConfig{
		DB: db,
	}
	// Go routine that runs separately from the main thread
	// This is a good place to put background tasks
	// go startScraping(db, 10, time.Minute)

	// Create a new router
	router := chi.NewRouter()

	// Set up CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Add Swagger UI route to the main router
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Test if the server is running
	router.Get("/ready", handlerReadiness)

	// Create V1 router
	v1 := chi.NewRouter()
	v1.Get("/err", handlerErr)
	v1.Post("/login", apiCfg.handlerLogin)
	v1.Post("/user", apiCfg.handlerCreateUser)
	v1.Get("/user", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1.Delete("/user", apiCfg.middlewareAuth(apiCfg.handlerDeleteUser))
	v1.Put("/user", apiCfg.middlewareAuth(apiCfg.handlerUpdateUser))

	v2 := chi.NewRouter()
	v2.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v2.Get("/feeds", apiCfg.middlewareAuth(apiCfg.handlerGetFeeds))
	v2.Put("/feeds/{feed_id}", apiCfg.middlewareAuth(apiCfg.handlerUpdateFeed))
	v2.Delete("/feeds/{feed_id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeed))

	v3 := chi.NewRouter()
	v3.Post("/follow", apiCfg.middlewareAuth(apiCfg.handlerFollowFeed))
	v3.Get("/follow", apiCfg.middlewareAuth(apiCfg.handlerGetFollows))
	v3.Delete("/follow/{feed_id}", apiCfg.middlewareAuth(apiCfg.handlerUnfollow))

	v4 := chi.NewRouter()
	v4.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPosts))
	// Mount the routers
	router.Mount("/v1", v1)
	router.Mount("/v2", v2)
	router.Mount("/v3", v3)
	router.Mount("/v4", v4)

	// Server configuration
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Printf("Server is running on port: %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
