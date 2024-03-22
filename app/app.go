package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/talut/dotenv"
	v1 "github.com/usercoredev/proto/api/v1"
	"github.com/usercoredev/usercore/app/services"
	"github.com/usercoredev/usercore/internal/cache"
	"github.com/usercoredev/usercore/internal/client"
	"github.com/usercoredev/usercore/internal/database"
	"github.com/usercoredev/usercore/internal/errorutil"
	"github.com/usercoredev/usercore/internal/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
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
			Scheme:             dotenv.GetString("TOKEN_SCHEME", "Bearer"),
			Issuer:             dotenv.MustGetString("APP_NAME"),
			Audience:           dotenv.MustGetString("JWT_AUDIENCE"),
			PrivateKeyPath:     dotenv.MustGetString("PRIVATE_KEY_PATH"),
			PublicKeyPath:      dotenv.MustGetString("PUBLIC_KEY_PATH"),
			AccessTokenExpire:  dotenv.GetDuration("ACCESS_TOKEN_EXPIRE", 1*time.Hour),
			RefreshTokenExpire: dotenv.GetDuration("REFRESH_TOKEN_EXPIRE", 24*time.Hour),
		},
		grpcServer: Server{
			Host: dotenv.GetString("GRPC_SERVER_HOST", ""),
			Port: dotenv.MustGetString("GRPC_SERVER_PORT"),
		},
		httpServer: Server{
			Host: dotenv.GetString("HTTP_SERVER_HOST", ""),
			Port: dotenv.MustGetString("HTTP_SERVER_PORT"),
		},
		databaseOptions: database.Database{
			Host:            dotenv.GetString("DB_HOST", ""),
			Port:            dotenv.GetString("DB_PORT", ""),
			User:            dotenv.GetString("DB_USER", ""),
			Password:        dotenv.GetString("DB_PASSWORD", ""),
			PasswordFile:    dotenv.GetString("DB_PASSWORD_FILE", ""),
			Database:        dotenv.GetString("DB_NAME", ""),
			DatabaseFile:    dotenv.GetString("DB_FILE_PATH", ""),
			Engine:          dotenv.GetString("DB_ENGINE", ""),
			Charset:         dotenv.GetString("DB_CHARSET", ""),
			Certificate:     dotenv.GetString("DB_CERTIFICATE_FILE", ""),
			EnableMigration: dotenv.GetString("DB_MIGRATE", ""),
		},
		cacheOptions: cache.Settings{
			Enabled:                    dotenv.GetBool("CACHE_ENABLED", false),
			Host:                       dotenv.MustGetString("CACHE_HOST"),
			Port:                       dotenv.MustGetString("CACHE_PORT"),
			Password:                   dotenv.GetString("CACHE_PASSWORD", ""),
			PasswordFile:               dotenv.GetString("CACHE_PASSWORD_FILE", ""),
			EncryptionKey:              dotenv.GetString("CACHE_ENCRYPTION_KEY", ""),
			UserCacheExpiration:        dotenv.GetDuration("USER_CACHE_EXPIRATION", 48*time.Hour),
			UserProfileCacheExpiration: dotenv.GetDuration("USER_PROFILE_CACHE_EXPIRATION", 48*time.Hour),
			UserCachePrefix:            dotenv.GetString("USER_CACHE_PREFIX", "user"),
			UserProfileCachePrefix:     dotenv.GetString("USER_PROFILE_CACHE_PREFIX", "profile"),
		},
		clientSettings: client.Settings{
			ClientFilePath: dotenv.MustGetString("CLIENTS_FILE_PATH"),
		},
	}
}

func (a *Application) ConnectToDatabase() {
	if err := a.databaseOptions.Connect(); err != nil {
		panic(err)
	}
}

func (a *Application) ConfigureToken() {
	a.tokenSettings.Setup()
}

func (a *Application) LoadClients() {
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
			client.ClientInterceptor(a.clientSettings),
			a.tokenSettings.AuthInterceptor(),
		),
	)
	a.registerGRPCServices(s)
	go func() {
		fmt.Println("GRPC Server running on: ", lis.Addr())
		if err := s.Serve(lis); err != nil {
			panic(errors.New(fmt.Sprintf("Code: %d, %s: %v", errorutil.ErrGRPCFailedToServe.Code, errorutil.ErrGRPCFailedToServe.Message, err)))
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
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(_ context.Context, req *http.Request) metadata.MD {
			return metadata.Pairs(string(client.Key), req.Header.Get(string(client.Key)))
		}),
	)
	a.registerHTTPServices(ctx, mux, conn)
	httpServerAddr := fmt.Sprintf("%s:%s", a.httpServer.Host, a.httpServer.Port)
	server := &http.Server{
		Addr:    httpServerAddr,
		Handler: mux,
	}
	fmt.Println("HTTP Server running on: ", httpServerAddr)
	if err = server.ListenAndServe(); err != nil {
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
