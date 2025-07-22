package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/alexliesenfeld/health/interceptors"
	"github.com/alexliesenfeld/health/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewChecker(db *gorm.DB, rdb *redis.Client) http.HandlerFunc {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	checker := health.NewChecker(
		health.WithTimeout(20*time.Second),
		health.WithInterceptors(interceptors.BasicLogger()),
		health.WithCacheDuration(0),

		health.WithCheck(health.Check{
			Name:    "Postgres Check",
			Timeout: 2 * time.Second,

			Check: func(ctx context.Context) error {
				var result int
				err := sqlDB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
				if err != nil {
					return err
				}
				return nil
			},

			Interceptors: []health.Interceptor{interceptors.BasicLogger()},

			MaxContiguousFails: 1,

			MaxTimeInError: 0,
		}),

		health.WithCheck(health.Check{
			Name: "Redis Check",

			Timeout: 2 * time.Second,

			Check: func(ctx context.Context) error {
				err := rdb.Ping(ctx).Err()
				if err != nil {
					return err
				}
				return nil
			},

			Interceptors: []health.Interceptor{interceptors.BasicLogger()},

			MaxContiguousFails: 1,

			MaxTimeInError: 0,
		}),
	)

	handler := health.NewHandler(checker,
		health.WithResultWriter(health.NewJSONResultWriter()),
		health.WithMiddleware(
			middleware.BasicLogger(),
			middleware.BasicAuth("user", "password"),
		),

		health.WithStatusCodeUp(http.StatusOK),
		health.WithStatusCodeDown(http.StatusServiceUnavailable),
	)

	return handler
}
