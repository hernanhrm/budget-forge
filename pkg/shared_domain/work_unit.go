package shared_domain

import "context"

type WorkUnit interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
