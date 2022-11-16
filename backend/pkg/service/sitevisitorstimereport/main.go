package sitevisitorstimereport

import (
	"github.com/poeticmetric/poeticmetric/backend/pkg/depot"
	"github.com/poeticmetric/poeticmetric/backend/pkg/service/sitereportfilters"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Datum struct {
	DateTime     time.Time `json:"dateTime"`
	VisitorCount uint64    `json:"visitorCount"`
}

type Report struct {
	AverageVisitorCount uint64    `json:"averageVisitorCount"`
	Data                []Datum   `json:"data"`
	Interval            *Interval `json:"interval"`
}

func Get(dp *depot.Depot, filters *sitereportfilters.Filters) (*Report, error) {
	interval := getInterval(filters)
	report := &Report{
		Interval: interval,
	}

	q := sitereportfilters.Apply(dp, filters)

	valueSubQuery := q.
		Select(
			strings.Join([]string{
				"toStartOfInterval(date_time, @timeWindowInterval, @timeZone) as date_time",
				"count(*) as visitor_count",
			}, ","),
			map[string]interface{}{
				"timeWindowInterval": gorm.Expr(interval.ToQuery()),
				"timeZone":           filters.GetTimeZone(),
			},
		).
		Group("date_time")

	fillerSubQuery := dp.ClickHouse().
		Raw(
			strings.Join([]string{
				"select",
				strings.Join([]string{
					"@start + interval arrayJoin(range(0, toUInt64(dateDiff('second', @start, @end)), @intervalSeconds)) second as date_time",
					"0 as visitor_count",
				}, ","),
			}, " "),
			map[string]any{
				"end":             filters.End,
				"intervalSeconds": interval.ToDuration().Seconds(),
				"start":           filters.Start,
			},
		)

	err := dp.ClickHouse().
		Table("((?) union all (?))", valueSubQuery, fillerSubQuery).
		Select(
			"date_time",
			"sum(visitor_count) as visitor_count",
		).
		Group("date_time").
		Order("date_time").
		Find(&report.Data).
		Error
	if err != nil {
		return nil, err
	}

	return report, nil
}
