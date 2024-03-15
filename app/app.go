package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/app/services"
	"github.com/usercoredev/usercore/cache"
	"github.com/usercoredev/usercore/database"
	"github.com/usercoredev/usercore/utils"
	"github.com/usercoredev/usercore/utils/client"
	"github.com/usercoredev/usercore/utils/server"
	"github.com/usercoredev/usercore/utils/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
)

type Application struct {
	clientSettings  client.Settings
	grpcServer      Server
	httpServer      Server
	tokenSettings   token.Settings
	databaseOptions database.Database
	cacheOptions    cache.Settings
}

type Server struct {
	Host string
	Port string
}

func Create() Application {
	return Application{
		tokenSettings: token.Settings{
			Scheme:             os.Getenv("TOKEN_SCHEME"),
			Issuer:             os.Getenv("APP_NAME"),
			Audience:           os.Getenv("JWT_AUDIENCE"),
			PrivateKeyPath:     os.Getenv("PRIVATE_KEY_PATH"),
			PublicKeyPath:      os.Getenv("PUBLIC_KEY_PATH"),
			AccessTokenExpire:  os.Getenv("ACCESS_TOKEN_EXPIRE"),
			RefreshTokenExpire: os.Getenv("REFRESH_TOKEN_EXPIRE"),
		},
		grpcServer: Server{
			Host: os.Getenv("GRPC_SERVER_HOST"),
			Port: os.Getenv("GRPC_SERVER_PORT"),
		},
		httpServer: Server{
			Host: os.Getenv("HTTP_SERVER_HOST"),
			Port: os.Getenv("HTTP_SERVER_PORT"),
		},
		databaseOptions: database.Database{
			Host:            os.Getenv("DB_HOST"),
			Port:            os.Getenv("DB_PORT"),
			User:            os.Getenv("DB_USER"),
			Password:        os.Getenv("DB_PASSWORD"),
			PasswordFile:    os.Getenv("DB_PASSWORD_FILE"),
			Database:        os.Getenv("DB_NAME"),
			DatabaseFile:    os.Getenv("DB_FILE_PATH"),
			Engine:          os.Getenv("DB_ENGINE"),
			Charset:         os.Getenv("DB_CHARSET"),
			Certificate:     os.Getenv("DB_CERTIFICATE_FILE"),
			EnableMigration: os.Getenv("DB_MIGRATE"),
		},
		cacheOptions: cache.Settings{
			Enabled:                    os.Getenv("CACHE_ENABLED"),
			Host:                       os.Getenv("CACHE_HOST"),
			Port:                       os.Getenv("CACHE_PORT"),
			Password:                   os.Getenv("CACHE_PASSWORD"),
			PasswordFile:               os.Getenv("CACHE_PASSWORD_FILE"),
			EncryptionKey:              os.Getenv("CACHE_ENCRYPTION_KEY"),
			UserCacheExpiration:        os.Getenv("USER_CACHE_EXPIRATION"),
			UserCachePrefix:            os.Getenv("USER_CACHE_PREFIX"),
			UserProfileCacheExpiration: os.Getenv("USER_PROFILE_CACHE_EXPIRATION"),
			UserProfileCachePrefix:     os.Getenv("USER_PROFILE_CACHE_PREFIX"),
		},
		clientSettings: client.Settings{
			ClientFilePath: os.Getenv("CLIENTS_FILE_PATH"),
		},
	}
}

func (a *Application) ConnectToDatabase() {
	if err := a.databaseOptions.Connect(); err != nil {
		panic(err)
	}
}

func (a *Application) ConfigureToken() {
	if a.tokenSettings.PrivateKeyPath == "" {
		panic(utils.ErrPrivateKeyPathNotSet)
	}
	if a.tokenSettings.PublicKeyPath == "" {
		panic(utils.ErrPublicKeyPathNotSet)
	}
	a.tokenSettings.Setup()
}

func (a *Application) LoadClients() {
	if a.clientSettings.ClientFilePath == "" {
		panic(utils.ErrClientFilePathNotSet)
	}
	if err := a.clientSettings.LoadClients(); err != nil {
		panic(err)
	}
}

func (a *Application) SetupCache() {
	if err := a.cacheOptions.SetupCache(); err != nil {
		panic(err)
	}
}

func (a *Application) StartServer() {
	if a.grpcServer.Port == "" {
		panic(utils.ErrGRPCPortNotSet)
	}
	address := fmt.Sprintf("%s:%s", a.grpcServer.Host, a.grpcServer.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	a.startGRPCServer(lis)
	err = a.startHTTPServer()
	if err != nil {
		panic(err)
	}
}

func (a *Application) startGRPCServer(lis net.Listener) {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			server.ClientInterceptor(a.clientSettings.Clients),
			server.AuthInterceptor,
		),
		grpc.ChainStreamInterceptor(),
	)
	a.registerGRPCServices(s)
	go func() {
		fmt.Println("GRPC Server running on: ", lis.Addr())
		if err := s.Serve(lis); err != nil {
			panic(errors.New(fmt.Sprintf("Code: %d, %s: %v", utils.ErrGRPCFailedToServe.Code, utils.ErrGRPCFailedToServe.Message, err)))
		}
	}()
}

func (a *Application) registerGRPCServices(server *grpc.Server) {
	v1.RegisterAuthenticationServiceServer(server, &services.AuthenticationServer{})
	v1.RegisterUserServiceServer(server, &services.UserServer{})
	v1.RegisterSessionServiceServer(server, &services.SessionServer{})
	v1.RegisterRoleServiceServer(server, &services.RoleServer{})
	v1.RegisterPermissionServiceServer(server, &services.PermissionServer{})
	reflection.Register(server)
}

func (a *Application) startHTTPServer() error {
	address := fmt.Sprintf("%s:%s", a.grpcServer.Host, a.grpcServer.Port)
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	gwMux := runtime.NewServeMux(
		runtime.WithMetadata(func(_ context.Context, req *http.Request) metadata.MD {
			return metadata.Pairs("client_id", req.Header.Get("Client_id"))
		}),
	)
	a.registerHTTPServices(ctx, gwMux, conn)
	httpServerAddr := fmt.Sprintf("%s:%s", a.httpServer.Host, a.httpServer.Port)
	gwServer := &http.Server{
		Addr:    httpServerAddr,
		Handler: gwMux,
	}
	fmt.Println("HTTP Server running on: ", httpServerAddr)
	if err = gwServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (a *Application) registerHTTPServices(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) {
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
