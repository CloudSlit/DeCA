// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"context"
	"github.com/flowshield/deca/internal/api"
	"github.com/flowshield/deca/internal/dao/certificate"
	"github.com/flowshield/deca/internal/initx"
	"github.com/flowshield/deca/internal/router"
	"github.com/flowshield/deca/internal/service"
)

// Injectors from wire.go:

func BuildInjector(ctx context.Context) (*Injector, func(), error) {
	execCloser, cleanup, err := initx.InitStorage()
	if err != nil {
		return nil, nil, err
	}
	ethClient, cleanup2, err := initx.InitEthClient(ctx)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	certificateRepo := &certificate.CertificateRepo{
		DB:  execCloser,
		Eth: ethClient,
	}
	cfsslHandler, err := initx.InitCfssl()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	tlsSrv := &service.TlsSrv{
		CertificateRepo: certificateRepo,
		CfsslHandler:    cfsslHandler,
	}
	tlsAPI := &api.TlsAPI{
		TlsSrv: tlsSrv,
	}
	cache := initx.InitOcspCache()
	ocspSrv := &service.OcspSrv{
		CertificateRepo: certificateRepo,
		CfsslHandler:    cfsslHandler,
		Cache:           cache,
		Ctx:             ctx,
	}
	ocspAPI := &api.OcspAPI{
		OcspSrv: ocspSrv,
	}
	certificateSrv := &service.CertificateSrv{
		CertificateRepo: certificateRepo,
	}
	certificateAPI := &api.CertificateAPI{
		CertificateSrv: certificateSrv,
	}
	routerRouter := &router.Router{
		TlsAPI:         tlsAPI,
		OcspAPI:        ocspAPI,
		CertificateAPI: certificateAPI,
	}
	engine := InitGinEngine(routerRouter)
	serveMux := InitOcspEngine(ocspAPI)
	injector := &Injector{
		Engine:     engine,
		OcspEngine: serveMux,
		Router:     routerRouter,
	}
	return injector, func() {
		cleanup2()
		cleanup()
	}, nil
}
