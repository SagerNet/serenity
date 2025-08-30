package server

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/sagernet/serenity/common/metadata"
	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/option"
	boxOption "github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func (s *Server) initializeRoutes() {
	s.chiRouter.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler)
	s.chiRouter.Get("/", s.render)
	s.chiRouter.Get("/{profileName}", s.render)
}

func (s *Server) render(writer http.ResponseWriter, request *http.Request) {
	profileName := chi.URLParam(request, "profileName")
	if profileName == "" {
		// compatibility with legacy versions
		profileName = request.URL.Query().Get("profile")
	}
	if strings.HasSuffix(profileName, "/") {
		profileName = profileName[:len(profileName)-1]
	}
	var profile *Profile
	if len(s.users) == 0 {
		if profileName == "" {
			profile = s.profile.DefaultProfile()
		} else {
			profile = s.profile.ProfileByName(profileName)
		}
	} else {
		user := s.authorization(request)
		if user == nil {
			writer.WriteHeader(http.StatusUnauthorized)
			s.accessLog(request, http.StatusUnauthorized, 0)
			return
		}
		if len(user.Profile) == 0 {
			writer.WriteHeader(http.StatusNotFound)
			s.accessLog(request, http.StatusNotFound, 0)
			return
		}
		if profileName == "" {
			profileName = user.DefaultProfile
		}
		if profileName == "" {
			profileName = user.Profile[0]
		}
		if !common.Contains(user.Profile, profileName) {
			writer.WriteHeader(http.StatusNotFound)
			s.accessLog(request, http.StatusNotFound, 0)
			return
		}
		profile = s.profile.ProfileByName(profileName)
	}
	if profile == nil {
		writer.WriteHeader(http.StatusNotFound)
		s.accessLog(request, http.StatusNotFound, 0)
		return
	}
	options, err := profile.Render(M.Detect(request.Header.Get("User-Agent")))
	if err != nil {
		s.logger.Error(E.Cause(err, "render options"))
		render.Status(request, http.StatusInternalServerError)
		render.PlainText(writer, request, err.Error())
		s.accessLog(request, http.StatusInternalServerError, len(err.Error()))
		return
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoderContext(s.ctx, &buffer)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(&options)
	if err != nil {
		s.logger.Error(E.Cause(err, "marshal options"))
		render.Status(request, http.StatusInternalServerError)
		render.PlainText(writer, request, err.Error())
		s.accessLog(request, http.StatusInternalServerError, len(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(buffer.Bytes())
	s.accessLog(request, http.StatusOK, buffer.Len())
}

func (s *Server) RenderHeadless(profileName string, metadata metadata.Metadata) (*boxOption.Options, error) {
	var profile *Profile
	if profileName == "" {
		s.profile.DefaultProfile()
	} else {
		profile = s.profile.ProfileByName(profileName)
	}
	if profile == nil {
		return nil, E.New("profile not found")
	}
	options, err := profile.Render(metadata)
	if err != nil {
		return nil, E.Cause(err, "render options")
	}
	return options, nil
}

func (s *Server) accessLog(request *http.Request, responseCode int, responseLen int) {
	var userString string
	if username, password, ok := request.BasicAuth(); ok {
		if responseCode == http.StatusUnauthorized {
			userString = username + ":" + password
		} else {
			userString = username
		}
	}
	s.logger.Debug("accepted ", request.RemoteAddr, " - ", userString, " \"", request.Method, " ", request.URL, " ", request.Proto, "\" ", responseCode, " ", responseLen, " \"", request.UserAgent(), "\"")
}

func (s *Server) authorization(request *http.Request) *option.User {
	username, password, ok := request.BasicAuth()
	if !ok {
		return nil
	}
	users, loaded := s.userMap[username]
	if !loaded {
		return nil
	}
	for _, user := range users {
		if user.Password == password {
			return &user
		}
	}
	return nil
}
