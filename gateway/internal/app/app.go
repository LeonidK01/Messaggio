package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LeonidK01/Messaggio/internal/delivery"
	"github.com/LeonidK01/Messaggio/internal/repository"
	"github.com/LeonidK01/Messaggio/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

func Start() error {
	appVersion := "v1"

	g := initGin()
	msgWriter := initKafkaWriter("kafka:9092", "messages")
	conn, err := initPostgres("postgres:5432")
	if err != nil {
		return fmt.Errorf("failed init PostgreSQL connection: %w", err)
	}

	mr := repository.NewMessagePostgresqlRepository(conn)
	mq := repository.NewMessageKafkaBroker(msgWriter)

	muc := usecase.NewMessageUsecase(mr, mq)

	app := g.Group(fmt.Sprintf("/%s", appVersion))

	delivery.HandleMessageGinDelivery(app.Group("/messages"), muc)

	srv := initHttp(":8080", g.Handler())
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed server shutdown: %w", err)
	}

	if err := msgWriter.Close(); err != nil {
		return fmt.Errorf("failed close kafka writer: %w", err)
	}

	if err := conn.Close(context.Background()); err != nil {
		return fmt.Errorf("failed close postgreSQL: %w", err)
	}

	<-ctx.Done()

	log.Println("Server exiting")

	return nil
}

func initGin() *gin.Engine {
	router := gin.Default()

	return router
}

func initHttp(addr string, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return srv
}

func initKafkaWriter(url string, topic string) *kafka.Writer {
	kw := &kafka.Writer{
		Addr:     kafka.TCP(url),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return kw
}

func initPostgres(url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	return conn, nil
}
