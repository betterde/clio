package auth

import (
	"context"
	"errors"
	"github.com/betterde/clio/global"
	"github.com/betterde/clio/internal/journal"
	"github.com/betterde/clio/internal/response"
	"github.com/gofiber/fiber/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

var (
	Authorizer   *authorization.Authorizer[*oauth.IntrospectionContext]
	Interception *Interceptor[*oauth.IntrospectionContext]
)

func init() {
	authorizer, err := authorization.New(global.Ctx, zitadel.New("zitadel.svc.dev"), oauth.DefaultAuthorization("authentication.json"))
	if err != nil {
		journal.Logger.Error(err)
	}

	Authorizer = authorizer
	Interception = New(Authorizer)
}

type Interceptor[T authorization.Ctx] struct {
	authorizer *authorization.Authorizer[T]
}

func New[T authorization.Ctx](authorizer *authorization.Authorizer[T]) *Interceptor[T] {
	return &Interceptor[T]{
		authorizer: authorizer,
	}
}

func (i *Interceptor[T]) RequireAuthorization(options ...authorization.CheckOption) func(next fiber.Handler) fiber.Handler {
	return func(next fiber.Handler) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			authCtx, err := i.authorizer.CheckAuthorization(ctx.Context(), ctx.Get(authorization.HeaderName), options...)
			if err != nil {
				if errors.Is(err, &authorization.UnauthorizedErr{}) {
					return ctx.Status(fiber.StatusUnauthorized).JSON(response.UnAuthenticated("Unauthorized"))
				}

				return ctx.Status(fiber.StatusForbidden).JSON(response.Forbidden("Forbidden"))
			}
			ctx.SetUserContext(authorization.WithAuthContext(ctx.Context(), authCtx))
			return next(ctx)
		}
	}
}

func (i *Interceptor[T]) Context(ctx context.Context) T {
	return authorization.Context[T](ctx)
}

func Middleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authCtx, err := Authorizer.CheckAuthorization(ctx.Context(), ctx.Get(authorization.HeaderName))
		if err != nil {
			if errors.Is(err, &authorization.UnauthorizedErr{}) {
				return ctx.Status(fiber.StatusUnauthorized).JSON(response.UnAuthenticated("Unauthorized"))
			}

			return ctx.Status(fiber.StatusForbidden).JSON(response.Forbidden("Forbidden"))
		}
		ctx.SetUserContext(authorization.WithAuthContext(ctx.Context(), authCtx))
		return ctx.Next()
	}
}
