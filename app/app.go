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
)

var clients []client.Client

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

type DatabaseOptions struct {
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

type Client struct {
	Clients    []client.Client
	ClientFile string
}

type CacheOptions struct {
	Enabled string
	Host    string
	Port    string
}

type App struct {
	DefaultServer
	Client          Client
	GRPCServer      Server
	HTTPServer      Server
	TokenOptions    TokenOptions
	DatabaseOptions DatabaseOptions
	CacheOptions    CacheOptions
}

func (app *App) LoadClients() {
	clientList, err := client.GetClients(app.Client.ClientFile)
	if err != nil {
		panic(err)
	}
	clients = append(app.Client.Clients, clientList...)
	app.Client.Clients = clients
}

func (app *App) ConnectToDatabase() {
	options := database.Database{
		Engine:       app.DatabaseOptions.Engine,
		DatabaseFile: app.DatabaseOptions.DatabaseFile,
		Host:         app.DatabaseOptions.Host,
		Port:         app.DatabaseOptions.Port,
		User:         app.DatabaseOptions.User,
		Password:     app.DatabaseOptions.Password,
		PasswordFile: app.DatabaseOptions.PasswordFile,
		Database:     app.DatabaseOptions.Database,
		Charset:      app.DatabaseOptions.Charset,
	}
	if err := options.Connect(); err != nil {
		panic(err)
	}
}

func (app *App) ConfigureToken() {
	token.SetPublicPrivateKey(app.TokenOptions.PublicKeyPath, app.TokenOptions.PrivateKeyPath)
	token.SetOptions(app.TokenOptions.AccessTokenExpire, app.TokenOptions.RefreshTokenExpire, app.TokenOptions.Scheme)
}

func (app *App) Cache() {
	if app.CacheOptions.Enabled != "true" {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic", r)
		}
	}()

	if err := cache.Redis(); err != nil {
		panic(err)
	}
}

func (app *App) StartServer() {
	if app.GRPCServer.Port == "" {
		panic("Port not set")
	}
	var lis net.Listener
	var err error
	address := fmt.Sprintf("%s:%s", app.GRPCServer.Host, app.GRPCServer.Port)

	if lis, err = net.Listen("tcp", address); err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(clientInterceptor, token.AuthInterceptor),
		grpc.ChainStreamInterceptor(),
	)
	app.RegisterGRPCServices(s)
	go func() {
		fmt.Println("GRPC Server running on: ", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()
	app.StartHTTPServer()
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

}

func (app *App) RegisterGRPCServices(server *grpc.Server) {
	v1.RegisterAuthenticationServiceServer(server, &services.AuthenticationServer{})
	v1.RegisterUserServiceServer(server, &services.UserServer{})
	v1.RegisterSessionServiceServer(server, &services.SessionServer{})
	v1.RegisterRoleServiceServer(server, &services.RoleServer{})
	v1.RegisterPermissionServiceServer(server, &services.PermissionServer{})
	reflection.Register(server)
}

func (app *App) StartHTTPServer() {
	address := fmt.Sprintf("%s:%s", app.GRPCServer.Host, app.GRPCServer.Port)
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
	app.RegisterHTTPServices(ctx, gwMux, conn)
	httpServerAddr := fmt.Sprintf("%s:%s", app.HTTPServer.Host, app.HTTPServer.Port)
	gwServer := &http.Server{
		Addr:    httpServerAddr,
		Handler: gwMux,
	}
	fmt.Println("HTTP Server running on: ", httpServerAddr)
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve gRPC-Gateway server: %v", err)
	}
}

func (app *App) RegisterHTTPServices(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) {
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
	mdClient := client.GetClient(clientID, clients)
	if mdClient == nil {
		return nil, status.Errorf(codes.Unauthenticated, responses.InvalidClient)
	}
	ctx = context.WithValue(ctx, "client", mdClient)
	return handler(ctx, req)
}
