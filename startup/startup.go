package startup

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nhutHao02/social-network-tweet-service/config"

	"github.com/nhutHao02/social-network-common-service/rabbitmq"
	"github.com/nhutHao02/social-network-tweet-service/database"
	"github.com/nhutHao02/social-network-tweet-service/internal"
	"github.com/nhutHao02/social-network-tweet-service/internal/api"
	"github.com/nhutHao02/social-network-tweet-service/pkg/redis"
	"github.com/nhutHao02/social-network-tweet-service/pkg/websocket"
	pb "github.com/nhutHao02/social-network-user-service/pkg/grpc"
	"golang.org/x/sync/errgroup"
)

func Start() {
	// init logger
	initLogger()
	// load congig
	cfg := config.LoadConfig()

	// run migration
	migration(cfg)

	// database setup
	db := database.OpenConnect(cfg.Database)

	// init redis
	rdb := redis.InitRedis(cfg.Redis)
	// Test connection
	_, err := rdb.Rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("failed to init redis------------", zap.Error(err))
		panic(fmt.Sprintf("Could not connect to Redis: %v", err))
	}

	// connect to grpc server
	userClient := openClientConnection(cfg.Client)

	// init Socket
	ws := websocket.NewSocket()

	// init RabbitMQ
	rabbitmq := initRabbitMQ(cfg.RabbitMQ)

	// init Server
	server := internal.InitializeServer(cfg, db, rdb, userClient, ws, rabbitmq)

	// run server
	runServer(server)

}

func runServer(server *api.Server) {
	var g errgroup.Group

	g.Go(func() error {
		return server.HTTPServer.RunHTTPServer()
	})

	// g.Go(func() error {
	// 	return server.GRPCServer.RunGRPCServer()
	// })

	if err := g.Wait(); err != nil {
		logger.Fatal("Error when start server", zap.Error(err))
	}
}

func initLogger() {
	err := logger.InitLogger()
	if err != nil {
		log.Fatalf("Could not initialize logger: %v", err)
	}
	defer logger.Sync()
}

func migration(cfg *config.Config) {
	// open connection to database
	db, err := sql.Open(cfg.Database.DbType, cfg.Database.ConnectionString)
	if err != nil {
		logger.Error("failed to connect to db when migration: ", zap.Error(err))
		return
	}
	// close connection when out
	defer db.Close()

	// create a migration instance
	m, err := migrate.New(
		strings.Join([]string{"file://", cfg.Database.MigrationFilePath}, ""),
		strings.Join([]string{cfg.Database.DbType, "://", cfg.Database.ConnectionString}, ""))
	if err != nil {
		logger.Error("failed to create a migration instance when migration: ", zap.Error(err))
		return
	}

	// migrate
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			logger.Error("failed to migrate UP file when migration: ", zap.Error(err))
		}
	}

	logger.Info("Migrations applied successfully")
}

func openClientConnection(cfg *config.ClientConfig) pb.UserServiceClient {
	conn, err := grpc.NewClient(cfg.UserService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to gRPC server: ", zap.Error(err))
	}
	client := pb.NewUserServiceClient(conn)
	logger.Info("Connect to user gRPC server port: " + cfg.UserService)
	return client
}

func initRabbitMQ(cfg *config.RabbitMQConfig) *rabbitmq.RabbitMQ {
	rabbit, err := rabbitmq.NewRabbitMQ(cfg.QueueName, cfg.ConnectionString)
	if err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
	}
	return rabbit
}
