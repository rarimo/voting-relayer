package service

import (
	"github.com/go-chi/chi"
	"github.com/rarimo/voting-relayer/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxRelayerConfig(s.cfg.RelayerConfig()),
		),
	)
	r.Route("/integrations/voting-relayer", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/vote", handlers.Voting)
		})
	})

	return r
}
