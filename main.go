package main

import (
	"database/sql"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	handler "restaurant-service/handler"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	restaurantHandler "restaurant-service/handler/restaurant"
	restaurantService "restaurant-service/service/restaurant"
	restaurantStore "restaurant-service/stores/restaurant"

	menuHandler "restaurant-service/handler/menu"
	menuService "restaurant-service/service/menu"
	menuStore "restaurant-service/stores/menu"
)

func main() {
	// Load environment variables from configs/.env file
	envPath := filepath.Join("configs", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s", envPath)
	}

	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// connection string
	connStr := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	migrationPath := "file://migrations"
	if err := runMigrations(migrationPath, connStr); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Kafka producer
	kafkaProducer, err := NewKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	log.Println("Migrations ran successfully")

	resStore := restaurantStore.New(db)
	resService := restaurantService.New(&resStore, kafkaProducer)
	resHandler := restaurantHandler.New(&resService)

	mStore := menuStore.New(db)
	mService := menuService.New(&mStore)
	mHandler := menuHandler.New(&mService)

	r := mux.NewRouter()

	// Define Restaurant
	r.HandleFunc("/restaurants", resHandler.Create).Methods("POST")
	r.HandleFunc("/restaurants", resHandler.Read).Methods("GET")
	r.HandleFunc("/restaurants/{id}", resHandler.Update).Methods("PATCH")
	r.HandleFunc("/restaurants/{id}", resHandler.Delete).Methods("DELETE")

	// Define Menu
	r.HandleFunc("/menu", mHandler.Create).Methods("POST")
	r.HandleFunc("/menu/{restaurant_id}", mHandler.GetMenu).Methods("GET")

	r.HandleFunc("/restaurant/{restaurant_id}/order/{order_id}", resHandler.UpdateOrderStatus).Methods("PUT")
	// r.HandleFunc("/menu", mHandler.Update).Methods("PATCH")
	// r.HandleFunc("/menu/{restaurant_id}", mHandler.Delete).Methods("DELETE")

	// Start HTTP server
	go func() {
		log.Println("Starting server on port 8001")
		if err := http.ListenAndServe(":8001", r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start Kafka consumer to listen for messages from the "order-placed" topic
	go func() {
		if err := startKafkaConsumer(&resService); err != nil {
			log.Fatalf("Failed to start Kafka consumer: %v", err)
		}
	}()

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	log.Println("Shutting down server...")
}

func runMigrations(migrationURL, databaseURL string) error {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationURL,
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func NewKafkaProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kafka producer: %v", err)
	}

	return producer, nil
}

func startKafkaConsumer(resService handler.Restaurant) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		return fmt.Errorf("error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("order-placed", 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("error subscribing to Kafka topic: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received message: %s\n", string(msg.Value))
			// Process the message
			if err := resService.ProcessOrderPlacedEvent(msg.Value); err != nil {
				log.Printf("Failed to process message: %v", err)
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %v", err)
		}
	}
}
