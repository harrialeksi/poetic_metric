package language

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/th0th/poeticmetric/backend/pkg/depot"
	"github.com/th0th/poeticmetric/backend/pkg/model"
	"github.com/th0th/poeticmetric/backend/pkg/pointer"
	"github.com/th0th/poeticmetric/backend/pkg/service/sitereport/filter"
	h "github.com/th0th/poeticmetric/backend/pkg/testhelper"
)

func TestGet(t *testing.T) {
	type TestDatum struct {
		Language     *string
		VisitorCount uint64
	}

	dp := h.NewDepot()

	_ = dp.WithPostgresTransaction(func(dp2 *depot.Depot) error {
		start, err := time.Parse("2006-01-02", "2022-01-01")
		assert.NoError(t, err)

		end, err := time.Parse("2006-01-02", "2022-12-31")
		assert.NoError(t, err)

		testData := []*TestDatum{
			{Language: nil, VisitorCount: uint64(gofakeit.IntRange(100, 200))},
		}

		var totalVisitorCount uint64

		for testDatumIndex := 0; testDatumIndex < 12; testDatumIndex += 1 {
			visitorCount := uint64(gofakeit.IntRange(1, 1000))
			totalVisitorCount += visitorCount

			testData = append(testData, &TestDatum{
				Language:     pointer.Get(fmt.Sprintf("language-%d", testDatumIndex+1)),
				VisitorCount: visitorCount,
			})
		}

		modelSite := h.Site(dp2, nil)
		modelEvents := []*model.Event{}

		for _, testDatum := range testData {
			var visitorIndex uint64

			for visitorIndex = 0; visitorIndex < testDatum.VisitorCount; visitorIndex += 1 {
				visitorId := uuid.NewString()
				viewCount := gofakeit.IntRange(1, 10)

				for viewIndex := 0; viewIndex < viewCount; viewIndex += 1 {
					modelEvents = append(modelEvents, &model.Event{
						DateTime:  gofakeit.DateRange(start, end),
						Id:        uuid.NewString(),
						Kind:      model.EventKindPageView,
						Language:  testDatum.Language,
						SiteId:    modelSite.Id,
						VisitorId: visitorId,
					})
				}
			}
		}

		filteredTestData := []*TestDatum{}

		for _, testDatum := range testData {
			if testDatum.Language != nil {
				filteredTestData = append(filteredTestData, testDatum)
			}
		}

		sort.Slice(filteredTestData, func(i, j int) bool {
			if filteredTestData[i].VisitorCount > filteredTestData[j].VisitorCount {
				return true
			}

			if filteredTestData[i].VisitorCount == filteredTestData[j].VisitorCount {
				return *filteredTestData[i].Language < *filteredTestData[j].Language
			}

			return false
		})

		expectedReport := &Report{
			PaginationCursor: &PaginationCursor{
				Language:     *filteredTestData[9].Language,
				VisitorCount: filteredTestData[9].VisitorCount,
			},
		}

		for _, testDatum := range filteredTestData[0:10] {
			expectedReport.Data = append(expectedReport.Data, &Datum{
				Language:          *testDatum.Language,
				VisitorCount:      testDatum.VisitorCount,
				VisitorPercentage: uint16(math.Round(100 * float64(testDatum.VisitorCount) / float64(totalVisitorCount))),
			})
		}

		err = dp2.ClickHouse().
			Create(&modelEvents).
			Error
		assert.NoError(t, err)

		report, err := Get(dp2, &filter.Filters{
			End:    end,
			SiteId: modelSite.Id,
			Start:  start,
		}, nil)
		assert.NoError(t, err)

		assert.Equal(t, expectedReport, report)

		err = dp2.ClickHouse().
			Exec("optimize table events_buffer").
			Error
		assert.NoError(t, err)

		err = dp2.ClickHouse().
			Table("events").
			Where("site_id = ?", modelSite.Id).
			Delete(nil).
			Error
		if err != nil {
			return err
		}

		return errors.New("")
	})

}
