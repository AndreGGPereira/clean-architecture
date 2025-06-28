package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"clean-architecture/configs"
	"clean-architecture/internal/event/handler"
	"clean-architecture/internal/infra/web/webserver"
	"clean-architecture/pkg/events"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(cfg.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(cfg.RMQUser, cfg.RMQPassword, cfg.RMQHost, cfg.RMQPort)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	//	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	//	listOrderUseCase := NewListOrderUseCase(db)

	webserver := webserver.NewWebServer(cfg.WebServerPort)
	//	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	//	webserver.AddHandler("/order", webOrderHandler.Create, "POST")
	//	webserver.AddHandler("/order", webOrderHandler.List, "GET")
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	//	orderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	//pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", cfg.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
	// 	CreateOrderUseCase: *createOrderUseCase,
	// 	ListOrderUseCase:   *listOrderUseCase,
	// }}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", cfg.GraphQLServerPort)
	http.ListenAndServe(":"+cfg.GraphQLServerPort, nil)
}

func getRabbitMQChannel(user, password, host, port string) *amqp.Channel {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	var conn *amqp.Connection
	var err error

	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(connectionString)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
