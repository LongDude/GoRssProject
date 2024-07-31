package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/LongDude/GoRssProject/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
  godotenv.Load(".env")

  portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port not specified; Error reading .env (PORT)")
	}

	dbURL := os.Getenv("GOOSE_DBSTRING")
	if dbURL == "" {
		log.Fatal("Couldnt acquire db_URL from .env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to DB:", err)
	}

	apiCfg := apiConfig{
		DB: db.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))

	// Создаём маршрутизатор для пути '/checkhealth'
	v1Router := chi.NewRouter()
	v1Router.Get("/checkhealth", handlerReadiness)
	v1Router.Get("/err", handler_err)
  // Access to users database
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPIKey))
  // Access to feeds database
  v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
  // Подписка на блоги
  v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
  v1Router.Get("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
  v1Router.Delete("/feed_follow/{feelFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

  // Присоединяем его как подмаршрут к основному
	// То есть после /v1 все пути обрабатывает v1Router
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server staring on port %v", portString)
	err = srv.ListenAndServe()

	// Произошла ошибка регистрации сервера
	if err != nil {
		log.Fatal(err)
	}
}
