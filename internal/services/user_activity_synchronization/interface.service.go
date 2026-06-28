package useractivitysynchronization

import (
	"context"

	"github.com/google/uuid"
)

type ServicePort interface {
	SyncUserActivityState(
		ctx context.Context,
		userID uuid.UUID,
		force bool,
	) error
}
