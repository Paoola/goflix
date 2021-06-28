package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type jsonMovie struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Duration    int    `json:"duration"`
	TrailerURL  string `json:"trailer_url"`
}

func (s *server) handleMovieList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies, err := s.store.GetMovies()
		if err != nil {
			log.Printf("Cannot load movies. err=%v", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = make([]jsonMovie, len(movies))

		for i, m := range movies {
			resp[i] = mapMovieToJson(m)
		}
		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleMovieDetail() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			log.Printf("Cannot parse id to int. err=%v", err)
			s.respond(rw, r, nil, http.StatusBadRequest)
		}

		m, err := s.store.GetMovieById(id)
		if err != nil {
			log.Printf("Cannot load movie. err=%v", err)
			s.respond(rw, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = mapMovieToJson(m)
		s.respond(rw, r, resp, http.StatusOK)

	}
}

func (s *server) handleMovieCreate() http.HandlerFunc {
	type request struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
		TrailerURL  string `json:"trailer_url"`
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(rw, r, &req)
		if err != nil {
			log.Printf("Cannot parse movie body. err=%v", err)
			s.respond(rw, r, nil, http.StatusBadRequest)
			return
		}

		m := &Movie{
			ID:          0,
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
			Duration:    req.Duration,
			TrailerURL:  req.TrailerURL,
		}

		err = s.store.CreateMovie(m)

		if err != nil {
			log.Printf("Cannot create movie in DB. err=%v", err)
			s.respond(rw, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = mapMovieToJson(m)
		s.respond(rw, r, resp, http.StatusOK)

	}
}

func mapMovieToJson(m *Movie) jsonMovie {
	return jsonMovie{
		ID:          m.ID,
		Title:       m.Title,
		ReleaseDate: m.ReleaseDate,
		Duration:    m.Duration,
		TrailerURL:  m.TrailerURL,
	}
}
