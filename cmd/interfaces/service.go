package interfaces

import (
	"context"
)

type TraderService interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
