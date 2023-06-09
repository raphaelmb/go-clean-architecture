package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"sync"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/raphaelmb/go-clean-architecture/configs"
	"github.com/raphaelmb/go-clean-architecture/internal/event/handler"
	"github.com/raphaelmb/go-clean-architecture/internal/infra/database"
	"github.com/raphaelmb/go-clean-architecture/internal/infra/graph"
	"github.com/raphaelmb/go-clean-architecture/internal/infra/grpc/pb"
	"github.com/raphaelmb/go-clean-architecture/internal/infra/grpc/service"
	"github.com/raphaelmb/go-clean-architecture/internal/infra/web/webserver"
	"github.com/raphaelmb/go-clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	database.TryCreateTable(db)

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrderUseCase(db)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		webserver := webserver.NewWebServer(configs.WebServerPort)
		webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
		webserver.AddHandler("/order", "post", webOrderHandler.Create)
		webserver.AddHandler("/order", "get", webOrderHandler.List)
		fmt.Println("Web server started on port", configs.WebServerPort)
		webserver.Start()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		grpcServer := grpc.NewServer()
		orderServices := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
		pb.RegisterOrderServiceServer(grpcServer, orderServices)
		reflection.Register(grpcServer)
		fmt.Println("GRPC server started on port", configs.GRPCServerPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
		if err != nil {
			panic(err)
		}
		grpcServer.Serve(lis)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			CreateOrderUseCase: *createOrderUseCase,
			ListOrdersUseCase:  *listOrdersUseCase,
		}}))
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", srv)
		fmt.Println("GraphQL server started on port", configs.GraphQLServerPort)
		http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
		wg.Done()
	}()

	wg.Wait()
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
