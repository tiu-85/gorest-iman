package adapters

import (
	"context"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/pkg/errors"
)

type queryStartTime struct {
	start time.Time
}
type LogHook struct {
	Logger Logger
}

func (d LogHook) BeforeQuery(ctx context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return context.WithValue(ctx, queryStartTimeKey, queryStartTime{start: time.Now()}), nil
}

func (d LogHook) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	queryTime := ctx.Value(queryStartTimeKey).(queryStartTime)
	logger := d.Logger.WithCtx(ctx, "duration", time.Since(queryTime.start).Milliseconds())
	if q.Err != nil {
		switch {
		case !errors.Is(q.Err, pg.ErrTxDone) && !errors.Is(q.Err, io.EOF):
			logger.Error(ctx, q.Err, q.Query)
		default:
			logger.Warn(q.Err)
		}
	}

	sql, err := q.FormattedQuery()
	if err != nil {
		logger.Error(err)
		return nil
	}

	re := regexp.MustCompile(`[\t\n]`)
	sql = strings.ReplaceAll(re.ReplaceAllString(sql, " "), "\"", "")
	logger.Info(sql)

	return nil
}
