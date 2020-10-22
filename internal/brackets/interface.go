package brackets

import "context"

type Brackets interface {
	Validate(ctx context.Context, str string) bool
	Fix(ctx context.Context, str string) (error, string)
}