package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/sagernet/serenity/common/cachefile"
	"github.com/sagernet/serenity/option"
	"github.com/sagernet/serenity/subscription"
	"github.com/sagernet/serenity/template"
	"github.com/sagernet/sing-box/common/tls"
	"github.com/sagernet/sing-box/log"
	boxOption "github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
	"github.com/sagernet/sing/service"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/http2"
)

type Server struct {
	createdAt    time.Time
	ctx          context.Context
	logFactory   log.Factory
	logger       log.Logger
	chiRouter    chi.Router
	httpServer   *http.Server
	tlsConfig    tls.ServerConfig
	cacheFile    *cachefile.CacheFile
	subscription *subscription.Manager
	template     *template.Manager
	profile      *ProfileManager
	users        []option.User
	userMap      map[string][]option.User
}

func New(ctx context.Context, options option.Options) (*Server, error) {
	ctx = service.ContextWithDefaultRegistry(ctx)
	createdAt := time.Now()
	logFactory, err := log.New(log.Options{
		Context:       ctx,
		Options:       common.PtrValueOrDefault(options.Log),
		DefaultWriter: os.Stderr,
		BaseTime:      createdAt,
	})
	if err != nil {
		return nil, E.Cause(err, "create log factory")
	}

	if err := logFactory.Start(); err != nil {
		return nil, E.Cause(err, "start log factory failed")
	}

	chiRouter := chi.NewRouter()
	httpServer := &http.Server{
		Addr:    options.Listen,
		Handler: chiRouter,
	}
	if httpServer.Addr == "" {
		if options.TLS != nil && options.TLS.Enabled {
			httpServer.Addr = ":443"
		} else {
			httpServer.Addr = ":80"
		}
	}
	var tlsConfig tls.ServerConfig
	if options.TLS != nil {
		tlsConfig, err = tls.NewServer(ctx, logFactory.NewLogger("tls"), common.PtrValueOrDefault(options.TLS))
		if err != nil {
			return nil, err
		}
	}
	var cacheFilePath string
	if options.CacheFile != "" {
		cacheFilePath = options.CacheFile
	} else {
		cacheFilePath = "cache.db"
	}
	cacheFile := cachefile.New(cacheFilePath)
	subscriptionManager, err := subscription.NewSubscriptionManager(
		ctx,
		logFactory.NewLogger("subscription"),
		cacheFile,
		options.Subscriptions)
	if err != nil {
		return nil, err
	}
	templateManager, err := template.NewManager(
		ctx,
		logFactory.NewLogger("template"),
		options.Templates)
	if err != nil {
		return nil, err
	}
	profileManager, err := NewProfileManager(
		ctx,
		logFactory.NewLogger("profile"),
		subscriptionManager,
		templateManager,
		common.Map(options.Outbounds, func(it boxOption.Listable[boxOption.Outbound]) []boxOption.Outbound {
			return it
		}),
		options.Profiles,
	)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string][]option.User)
	for _, user := range options.Users {
		userMap[user.Name] = append(userMap[user.Name], user)
	}
	return &Server{
		createdAt:    createdAt,
		ctx:          ctx,
		logFactory:   logFactory,
		logger:       logFactory.Logger(),
		chiRouter:    chiRouter,
		httpServer:   httpServer,
		tlsConfig:    tlsConfig,
		cacheFile:    cacheFile,
		subscription: subscriptionManager,
		template:     templateManager,
		profile:      profileManager,
		users:        options.Users,
		userMap:      userMap,
	}, nil
}

func (s *Server) Start() error {
	s.initializeRoutes()
	err := s.cacheFile.Start()
	if err != nil {
		return err
	}
	err = s.subscription.Start()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}
	if s.tlsConfig != nil {
		err = s.tlsConfig.Start()
		if err != nil {
			return err
		}
		err = http2.ConfigureServer(s.httpServer, new(http2.Server))
		if err != nil {
			return err
		}
		stdConfig, err := s.tlsConfig.Config()
		if err != nil {
			return err
		}
		s.httpServer.TLSConfig = stdConfig
	}
	s.logger.Info("server started at ", listener.Addr())
	go func() {
		if s.httpServer.TLSConfig != nil {
			err = s.httpServer.ServeTLS(listener, "", "")
		} else {
			err = s.httpServer.Serve(listener)
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server serve error: ", err)
		}
	}()
	err = s.postStart()
	if err != nil {
		return err
	}
	s.logger.Info("serenity started (", F.Seconds(time.Since(s.createdAt).Seconds()), "s)")
	return nil
}

func (s *Server) postStart() error {
	err := s.subscription.PostStart()
	if err != nil {
		return E.Cause(err, "post-start subscription manager")
	}
	return nil
}

func (s *Server) Close() error {
	return common.Close(
		s.logFactory,
		common.PtrOrNil(s.httpServer),
		s.tlsConfig,
		common.PtrOrNil(s.cacheFile),
	)
}
