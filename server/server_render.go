package server

import (
	"bytes"
	"encoding/json"
	"net/http"

	M "github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"

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
			return
		}
		if len(user.Profile) == 0 {
			writer.WriteHeader(http.StatusNotFound)
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
			return
		}
		profile = s.profile.ProfileByName(profileName)
	}
	if profile == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	metadata := M.Detect(request.Header.Get("User-Agent"))
	options, err := profile.Render(metadata)
	if err != nil {
		s.logger.Error(E.Cause(err, "render options"))
		render.Status(request, http.StatusInternalServerError)
		render.PlainText(writer, request, err.Error())
		return
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(&options)
	if err != nil {
		s.logger.Error(E.Cause(err, "marshal options"))
		render.Status(request, http.StatusInternalServerError)
		render.PlainText(writer, request, err.Error())
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(buffer.Bytes())
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
