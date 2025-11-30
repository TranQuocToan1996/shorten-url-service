package webhook

import "context"

type Client interface {
	Notify(ctx context.Context, callbackURL string, payload any) error
}
