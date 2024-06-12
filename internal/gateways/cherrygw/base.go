package cherrygw

import (
	"fmt"

	"github.com/bartosian/notibot/internal/core/ports"
	"github.com/bartosian/notibot/pkg/l0g"
	"github.com/cherryservers/cherrygo"
)

type Gateway struct {
	client *cherrygo.Client
	logger l0g.Logger
}

func NewCherryGateway(logger l0g.Logger) (ports.CherryGateway, error) {
	client, err := cherrygo.NewClient()
	if err != nil {
		err = fmt.Errorf("error creating cherry client: %w", err)
		logger.Error("error creating cherry client:", err, nil)

		return nil, err
	}

	return &Gateway{
		client: client,
		logger: logger,
	}, nil
}
