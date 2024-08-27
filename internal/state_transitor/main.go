package ingester

import (
	"context"
	"fmt"
	"sync"
	"time"

	rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"github.com/rarimo/voting-relayer/internal/config"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const (
	stateQuery = "operation_signed.operation_type='PASSPORT_ROOT_UPDATE'"
)

type Service struct {
	Processor
	log             *logan.Entry
	client          *http.HTTP
	catchupDisabled bool
	ws              *sync.WaitGroup
}

type Processor interface {
	Catchup(ctx context.Context) error
	Process(ctx context.Context, confirmationID string) error
	Query() string
	Name() string
}

func NewService(cfg config.Config, processor Processor, ws *sync.WaitGroup) *Service {
	return &Service{
		Processor: processor,
		log:       cfg.Log(),
		client:    cfg.Tendermint(),
		ws:        ws,
	}
}

func (s *Service) Run(ctx context.Context) {
	if s.ws != nil {
		defer s.ws.Done()
	}

	if !s.catchupDisabled {
		if err := s.Catchup(ctx); err != nil {
			s.log.WithError(err).Error("failed to catchup")
		}
	}

	running.WithBackOff(
		ctx, s.log, s.Processor.Name(), s.run,
		5*time.Second, 5*time.Second, 5*time.Second,
	)
}

func (s *Service) run(ctx context.Context) error {
	s.log.Info("Starting subscription")
	defer s.log.Info("Subscription finished")

	const depositChanSize = 100

	out, err := s.client.Subscribe(
		ctx,
		s.Processor.Name(),
		s.Processor.Query(),
		depositChanSize,
	)

	if err != nil {
		return errors.Wrap(err, "can not subscribe")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c := <-out:
			confirmation := c.Events[fmt.Sprintf("%s.%s", rarimocore.EventTypeOperationSigned, rarimocore.AttributeKeyConfirmationId)][0]
			s.log.Infof("New confirmation found %s", confirmation)

			if err := s.Process(ctx, confirmation); err != nil {
				s.log.WithError(err).Error("failed to process confirmation")
			}
		}
	}
}
