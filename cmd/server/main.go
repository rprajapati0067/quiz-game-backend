package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"google.golang.org/grpc"

	"github.com/rprajapati0067/quiz-app-tools/logger"
	"github.com/rprajapati0067/quiz-game-backend/initilization"
	authrpc "github.com/rprajapati0067/quiz-game-backend/rpc/auth"
	questionrpc "github.com/rprajapati0067/quiz-game-backend/rpc/question"
	userrpc "github.com/rprajapati0067/quiz-game-backend/rpc/user"

	"github.com/rprajapati0067/quiz-game-backend/internal/handlers"
	"github.com/rprajapati0067/quiz-game-backend/internal/repository"
	"github.com/rprajapati0067/quiz-game-backend/internal/service"
)

var httpAdapter *httpadapter.HandlerAdapter

func initServices() (service.AuthService, service.UserService, service.QuestionService) {
	userRepo := repository.NewMemoryUserRepository()
	questionRepo := repository.NewMemoryQuestionRepository()

	authSvc := service.NewAuthService(userRepo)
	userSvc := service.NewUserService(userRepo)
	questionSvc := service.NewQuestionService(questionRepo)

	return authSvc, userSvc, questionSvc
}

func setupHTTPMux() *http.ServeMux {
	authSvc, userSvc, questionSvc := initServices()

	httpHandlers := handlers.NewHTTPHandlers(authSvc, userSvc, questionSvc)
	mux := http.NewServeMux()
	httpHandlers.SetupRoutes(mux)

	return mux
}

func setupLocalServer() {
	logger.Info("Starting in LOCAL mode")

	authSvc, userSvc, questionSvc := initServices()

	// Setup HTTP REST API server
	httpMux := setupHTTPMux()
	go func() {
		logger.Info("HTTP REST API server running on :8080")
		logger.Info("API endpoints available at http://localhost:8080/api/v1/")
		if err := http.ListenAndServe(":8080", httpMux); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Setup gRPC server
	authHandler := handlers.NewAuthHandler(authSvc)
	userHandler := handlers.NewUserHandler(userSvc)
	questionHandler := handlers.NewQuestionHandler(questionSvc)

	grpcServer := grpc.NewServer()

	authrpc.RegisterAuthServiceServer(grpcServer, authHandler)
	userrpc.RegisterUserServiceServer(grpcServer, userHandler)
	questionrpc.RegisterQuestionServiceServer(grpcServer, questionHandler)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logger.Info("gRPC server running on :8082")
	logger.Info("Health check available at http://localhost:8080/health")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func setupLambdaServer() {
	logger.Info("Starting in LAMBDA mode")

	mux := setupHTTPMux()
	httpAdapter = httpadapter.New(mux)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return httpAdapter.ProxyWithContext(ctx, req)
	})
}

func main() {
	// Initialize logger
	initilization.Init()

	// Check if running in Lambda environment
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		setupLambdaServer()
	} else {
		setupLocalServer()
	}
}
