package serenity

import (
	"context"
	"encoding/json"
	"errors"
	std_log "log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/tls"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/http2"
)

var _ adapter.Service = (*Server)(nil)

type Server struct {
	boxService         *box.Box
	logger             logger.Logger
	logFactory         log.Factory
	httpServer         *http.Server
	tlsConfig          tls.ServerConfig
	subscriptionClient *SubscriptionClient
	profiles           []*Profile
	outbounds          []option.Outbound
	defaultProfile     string
}

func NewServer(ctx context.Context, options Options) (*Server, error) {
	createdAt := time.Now()
	boxService := common.Must1(box.New(box.Options{
		Options: option.Options{
			Log: &option.LogOptions{
				Disabled: true,
			},
		},
	}))
	common.Must(boxService.Start())

	logFactory, err := log.New(log.Options{
		Options:       common.PtrValueOrDefault(options.Log),
		DefaultWriter: os.Stderr,
		BaseTime:      createdAt,
	})
	if err != nil {
		return nil, E.Cause(err, "create log factory")
	}

	chiRouter := chi.NewRouter()
	server := &Server{
		boxService: boxService,
		logger:     logFactory.Logger(),
		logFactory: logFactory,
		httpServer: &http.Server{
			Handler:  chiRouter,
			ErrorLog: std_log.New(os.Stderr, "http", 0),
		},
		defaultProfile: options.DefaultProfile,
		outbounds:      options.Outbounds,
	}
	listen := options.Listen
	if listen == "" {
		if options.TLS != nil && options.TLS.Enabled {
			listen = ":443"
		} else {
			listen = ":80"
		}
	}
	server.httpServer.Addr = listen
	if options.TLS != nil {
		tlsConfig, err := tls.NewServer(ctx, logFactory.NewLogger("tls"), common.PtrValueOrDefault(options.TLS))
		if err != nil {
			return nil, err
		}
		server.tlsConfig = tlsConfig
	}
	server.subscriptionClient = NewSubscriptionClient(server.logger, options.Subscriptions)
	for i, profileOptions := range options.Profiles {
		profile, err := NewProfile(profileOptions)
		if err != nil {
			return nil, E.Cause(err, "parse profile[", i, "]")
		}
		server.profiles = append(server.profiles, profile)
	}
	if len(server.profiles) == 0 {
		return nil, errors.New("empty profiles")
	}
	chiRouter.Get("/", server.ServeHTTP)
	return server, nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	profileName := request.URL.Query().Get("profile")
	if profileName == "" {
		profileName = s.defaultProfile
	}
	if profileName == "" {
		profileName = s.profiles[0].name
	}
	profile := common.Find(s.profiles, func(it *Profile) bool {
		return it.Name() == profileName
	})
	if profile == nil {
		s.logger.Warn("profile not found: ", profileName)
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if !profile.CheckBasicAuthorization(request) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	var platform string
	userAgent := request.Header.Get("User-Agent")
	if strings.HasPrefix(userAgent, "SFA") {
		platform = PlatformAndroid
	} else if strings.HasPrefix(userAgent, "SFI") {
		platform = PlatformiOS
	} else if strings.HasPrefix(userAgent, "SFM") {
		platform = PlatformMacOS
	} else if strings.HasPrefix(userAgent, "SFT") {
		platform = PlatformAppleTVOS
	}
	var versionName string
	if strings.Contains(userAgent, "sing-box ") {
		versionName = strings.Split(userAgent, "sing-box ")[1]
		versionName = strings.Split(versionName, " ")[0]
		versionName = strings.Split(versionName, ")")[0]
	}
	var versionPtr *Version
	if versionName != "" {
		version := ParseVersion(versionName)
		versionPtr = &version
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	subscriptions := s.subscriptionClient.Subscriptions()
	options := profile.GenerateConfig(platform, versionPtr, s.outbounds, subscriptions)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(options)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Server) Start() error {
	err := s.subscriptionClient.Start()
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
		s.logger.Error("!!!")
	}()
	return nil
}

func (s *Server) Close() error {
	return common.Close(
		common.PtrOrNil(s.httpServer),
		s.tlsConfig,
		s.subscriptionClient,
		s.logFactory,
	)
}
