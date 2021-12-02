package apiserver

import (
	"gobackend/pkg/log"
	"gobackend/pkg/shutdown"
	"gobackend/pkg/shutdown/shutdownmanagers/posixsignal"

	"gobackend/internal/app/apiserver/config"
	"gobackend/internal/app/apiserver/store/mysql"
	genericoptions "gobackend/internal/pkg/options"
	genericserver "gobackend/internal/pkg/server"
)

type apiServer struct {
	gs               *shutdown.GracefulShutdown
	genericAPIServer *genericserver.GenericAPIServer
	mysqlOptions     *genericoptions.MySQLOptions
}

type preparedAPIServer struct {
	*apiServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().NewServer()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gs:               gs,
		genericAPIServer: genericServer,
		mysqlOptions:     cfg.MySQL,
	}

	return server, nil
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	if err := mysql.InitMySQLFactory(s.mysqlOptions); err != nil {
		log.Fatalf("init mysql failed: %s", err)
	}

	initRouter(s.genericAPIServer.Engine)

	s.gs.AddShutdownCallback(shutdown.Func(func(string) error {
		mysqlStore := mysql.GetMysqlFactory()
		if mysqlStore != nil {
			return mysqlStore.Close()
		}

		s.genericAPIServer.Close()

		return nil
	}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.genericAPIServer.Run()
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericserver.Config, lastErr error) {
	genericConfig = genericserver.NewConfig()

	lastErr = cfg.ApplyTo(genericConfig)

	return
}
