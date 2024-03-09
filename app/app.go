package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/app/responses"
	"github.com/usercoredev/usercore/app/services"
	"github.com/usercoredev/usercore/cache"
	"github.com/usercoredev/usercore/database"
	"github.com/usercoredev/usercore/utils/client"
	"github.com/usercoredev/usercore/utils/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"os"
)

var App Application

type Application struct {
	DefaultServer
	clientSettings  clientSettings
	grpcServer      Server
	httpServer      Server
	tokenOptions    TokenOptions
	databaseOptions databaseOptions
	cacheOptions    cacheOptions
}

type DefaultServer interface {
	StartServer()
	StartHTTPServer()
	RegisterHTTPServices()
	RegisterGRPCServices(server *grpc.Server)
	LoadClients()
	ConnectToDatabase()
	ConnectToCache()
	ConfigureToken()
}

type Server struct {
	Host string
	Port string
}

type TokenOptions struct {
	Scheme             string
	PrivateKeyPath     string
	PublicKeyPath      string
	AccessTokenExpire  string
	RefreshTokenExpire string
}

type databaseOptions struct {
	Engine       string
	DatabaseFile string
	Host         string
	Port         string
	User         string
	Password     string
	PasswordFile string
	Database     string
	Charset      string
}

type clientSettings struct {
	clients    []client.Client
	clientFile string
}

type cacheOptions struct {
	Enabled string
	Host    string
	Port    string
}

func Create() {
	App = Application{
		tokenOptions: TokenOptions{
			Scheme:             os.Getenv("TOKEN_SCHEME"),
			PrivateKeyPath:     os.Getenv("PRIVATE_KEY_PATH"),
			PublicKeyPath:      os.Getenv("PUBLIC_KEY_PATH"),
			AccessTokenExpire:  os.Getenv("ACCESS_TOKEN_EXPIRE"),
			RefreshTokenExpire: os.Getenv("REFRESH_TOKEN_EXPIRE"),
		},
		grpcServer: Server{
			Port: os.Getenv("GRPC_SERVER_PORT"),
		},
		httpServer: Server{
			Port: os.Getenv("HTTP_SERVER_PORT"),
		},
		databaseOptions: databaseOptions{
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			User:         os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			PasswordFile: os.Getenv("DB_PASSWORD_FILE"),
			Database:     os.Getenv("DB_NAME"),
			DatabaseFile: os.Getenv("DB_FILE_PATH"),
			Engine:       os.Getenv("DB_ENGINE"),
		},
		cacheOptions: cacheOptions{
			Enabled: os.Getenv("CACHE_ENABLED"),
			Host:    os.Getenv("CACHE_HOST"),
			Port:    os.Getenv("CACHE_PORT"),
		},
		clientSettings: clientSettings{
			clientFile: os.Getenv("CLIENTS_FILE_PATH"),
		},
	}
}

func (a *Application) LoadClients() {
	clientList, err := client.GetClients(a.clientSettings.clientFile)
	if err != nil {
		panic(err)
	}
	if len(clientList) == 0 {
		panic("No clients found")
	}
	a.clientSettings.clients = clientList
}

func (a *Application) ConnectToDatabase() {
	options := database.Database{
		Engine:       a.databaseOptions.Engine,
		DatabaseFile: a.databaseOptions.DatabaseFile,
		Host:         a.databaseOptions.Host,
		Port:         a.databaseOptions.Port,
		User:         a.databaseOptions.User,
		Password:     a.databaseOptions.Password,
		PasswordFile: a.databaseOptions.PasswordFile,
		Database:     a.databaseOptions.Database,
		Charset:      a.databaseOptions.Charset,
	}
	if err := options.Connect(); err != nil {
		panic(err)
	}
}

func (a *Application) ConfigureToken() {
	token.SetPublicPrivateKey(a.tokenOptions.PublicKeyPath, a.tokenOptions.PrivateKeyPath)
	token.SetOptions(a.tokenOptions.AccessTokenExpire, a.tokenOptions.RefreshTokenExpire, a.tokenOptions.Scheme)
}

func (a *Application) Cache() {
	if a.cacheOptions.Enabled != "true" {
		return
	}

	if err := cache.Redis(); err != nil {
		fmt.Println("Redis connection failed:", err)
	}
}

func (a *Application) StartServer() {
	if a.grpcServer.Port == "" {
		panic("gRPC Server port not set")
	}
	var lis net.Listener
	var err error
	address := fmt.Sprintf("%s:%s", a.grpcServer.Host, a.grpcServer.Port)

	if lis, err = net.Listen("tcp", address); err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(clientInterceptor, token.AuthInterceptor),
		grpc.ChainStreamInterceptor(),
	)
	a.RegisterGRPCServices(s)
	go func() {
		fmt.Println("GRPC Server running on: ", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()
	a.StartHTTPServer()
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

}

func (a *Application) RegisterGRPCServices(server *grpc.Server) {
	v1.RegisterAuthenticationServiceServer(server, &services.AuthenticationServer{})
	v1.RegisterUserServiceServer(server, &services.UserServer{})
	v1.RegisterSessionServiceServer(server, &services.SessionServer{})
	v1.RegisterRoleServiceServer(server, &services.RoleServer{})
	v1.RegisterPermissionServiceServer(server, &services.PermissionServer{})
	reflection.Register(server)
}

func (a *Application) StartHTTPServer() {
	address := fmt.Sprintf("%s:%s", a.grpcServer.Host, a.grpcServer.Port)
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwMux := runtime.NewServeMux(
		runtime.WithMetadata(func(_ context.Context, req *http.Request) metadata.MD {
			return metadata.Pairs("client_id", req.Header.Get("Client_id"))
		}),
	)
	a.RegisterHTTPServices(ctx, gwMux, conn)
	httpServerAddr := fmt.Sprintf("%s:%s", a.httpServer.Host, a.httpServer.Port)
	gwServer := &http.Server{
		Addr:    httpServerAddr,
		Handler: gwMux,
	}
	fmt.Println("HTTP Server running on: ", httpServerAddr)
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve gRPC-Gateway server: %v", err)
	}
}

func (a *Application) RegisterHTTPServices(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) {
	if err := v1.RegisterAuthenticationServiceHandler(ctx, mux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	if err := v1.RegisterUserServiceHandler(ctx, mux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	if err := v1.RegisterSessionServiceHandler(ctx, mux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	if err := v1.RegisterRoleServiceHandler(ctx, mux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	if err := v1.RegisterPermissionServiceHandler(ctx, mux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
}

func clientInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	if len(md.Get("client_id")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, responses.ClientRequired)
	}
	clientID := md.Get("client_id")[0]
	if clientID == "" {
		return nil, status.Errorf(codes.Unauthenticated, responses.ClientRequired)
	}
	mdClient := client.GetClient(clientID, App.clientSettings.clients)
	if mdClient == nil {
		return nil, status.Errorf(codes.Unauthenticated, responses.InvalidClient)
	}
	ctx = context.WithValue(ctx, "client", mdClient)
	return handler(ctx, req)
}
