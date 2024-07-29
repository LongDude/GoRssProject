package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
  "github.com/joho/godotenv"
)

func main(){
  godotenv.Load(".env")
  
  portString := os.Getenv("PORT")
  if portString == ""{
    log.Fatal("Port not specified; Error reading .env (PORT)")
  }
  
  router := chi.NewRouter()

  router.Use(cors.Handler(cors.Options{
    AllowedOrigins: []string{"https://*", "http://*"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowCredentials: false,
    AllowedHeaders: []string{"*"},
    ExposedHeaders: []string{"Link"},                   
    MaxAge: 300,
  }))

  // Создаём маршрутизатор для пути '/checkhealth'
  v1Router := chi.NewRouter()
  v1Router.Get("/checkhealth", handlerReadiness)
  v1Router.Get("/err", handler_err)
  // Присоединяем его как подмаршрут к основному
  // То есть после /v1 все пути обрабатывает v1Router
  router.Mount("/v1", v1Router)

  srv := &http.Server{
    Handler: router,
    Addr: ":" + portString,
  }

  log.Printf("Server staring on port %v", portString)
  err := srv.ListenAndServe()

  // Произошла ошибка регистрации сервера
  if err != nil{
    log.Fatal(err)
  }
}
