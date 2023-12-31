package app

import (
	"context"
	"github.com/Nol1feee/CLI-chat/auth/pkg/closer"

	//"context"
	"github.com/Nol1feee/CLI-chat/auth/internal/api/auth"
	"github.com/Nol1feee/CLI-chat/auth/internal/config"
	"github.com/Nol1feee/CLI-chat/auth/internal/repository"
	authRepo "github.com/Nol1feee/CLI-chat/auth/internal/repository/auth"
	"github.com/Nol1feee/CLI-chat/auth/internal/service"
	authService "github.com/Nol1feee/CLI-chat/auth/internal/service/auth"
	"github.com/jackc/pgx/v4/pgxpool"
	//"google.golang.org/genproto/googleapis/appengine/v1"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	httpConfig config.HTTPConfig

	pool *pgxpool.Pool

	authService    service.AuthService
	authRepository repository.AuthRepository

	authImplementation *auth.Implementation

	closer *closer.Closer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetSwaggerConfig() config.SwaggerCfg {
	cfg, err := config.NewSwaggerConfig()
	if err != nil {
		log.Fatalf("NewSwaggerConfig | %s\n", err)
	}
	return cfg
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {
	pgCfg, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("NewPGConfig | %s\n", err)
	}
	return pgCfg
}

func (s *serviceProvider) GetGRPCConfig() config.GRPCConfig {
	cfgGRPC, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("NewGRPCConfig | %s\n", err)
	}

	return cfgGRPC
}

func (s *serviceProvider) GetHTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfgHTTP, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("NewHTTPConfig | %s\n", err)
		}
		s.httpConfig = cfgHTTP
	}

	return s.httpConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pool == nil {
		con, err := pgxpool.Connect(ctx, s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("connect to DB error | %s ", err)
		}

		err = con.Ping(ctx)
		if err != nil {
			log.Fatalf("ping DB error | %s ", err)
		}

		s.closer.Add(func() error {
			con.Close()
			return nil
		})

		s.pool = con
	}

	return s.pool
}

func (s *serviceProvider) AuthRepo(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepo.NewRepo(s.PgPool(ctx))
	}
	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepo(ctx))
	}
	return s.authService
}

func (s *serviceProvider) NewImplementation(ctx context.Context) *auth.Implementation {
	if s.authImplementation == nil {
		s.authImplementation = auth.NewImplementation(s.AuthService(ctx))
	}
	return s.authImplementation
}
