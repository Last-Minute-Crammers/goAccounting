package transcationService

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/global/db"
	transactionModel "goAccounting/internal/model/transaction"
	"goAccounting/util/timeTool"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type StatisticService struct{}

func NewStatisticService() *StatisticService {
	return &StatisticService{}
}

// PeriodType represents different time periods
type PeriodType string

const (
	Daily   PeriodType = "daily"
	Weekly  PeriodType = "weekly"
	Monthly PeriodType = "monthly"
	Yearly  PeriodType = "yearly"
)

// PeriodStatistic represents aggregated statistics for a period
type PeriodStatistic struct {
	Period     string             `json:"period"`
	StartTime  time.Time          `json:"start_time"`
	EndTime    time.Time          `json:"end_time"`
	Statistics global.IEStatistic `json:"statistics"`
}

// GetPeriodStatistics aggregates daily statistics into period-based statistics
func (s *StatisticService) GetPeriodStatistics(
	userId uint,
	periodType PeriodType,
	startTime, endTime time.Time,
	categoryIds []uint,
	ie *constant.IncomeExpense,
	ctx context.Context,
) ([]PeriodStatistic, error) {
	periods := s.calculatePeriods(periodType, startTime, endTime)
	results := make([]PeriodStatistic, 0, len(periods))

	dao := transactionModel.NewStatisticDao(db.GetDb(ctx))

	for _, period := range periods {
		condition := s.buildStatisticCondition(userId, period.Start, period.End, categoryIds)
		stats, err := dao.GetIeStatisticByCondition(ie, *condition)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get period statistics")
		}

		results = append(results, PeriodStatistic{
			Period:     period.Label,
			StartTime:  period.Start,
			EndTime:    period.End,
			Statistics: stats,
		})
	}

	return results, nil
}

// GetCategoryPeriodStatistics gets category-level statistics for a specific period
func (s *StatisticService) GetCategoryPeriodStatistics(
	userId uint,
	categoryIds []uint,
	periodType PeriodType,
	startTime, endTime time.Time,
	ie constant.IncomeExpense,
	ctx context.Context,
) (global.AmountCount, error) {
	condition := s.buildStatisticCondition(userId, startTime, endTime, categoryIds)
	dao := transactionModel.NewStatisticDao(db.GetDb(ctx))

	return dao.GetAmountCountByCondition(*condition, ie)
}

type periodRange struct {
	Label string
	Start time.Time
	End   time.Time
}

func (s *StatisticService) calculatePeriods(periodType PeriodType, startTime, endTime time.Time) []periodRange {
	var periods []periodRange
	current := timeTool.ToDay(startTime)
	end := timeTool.ToDay(endTime)

	switch periodType {
	case Weekly:
		for current.Before(end) || current.Equal(end) {
			weekStart := s.getWeekStart(current)
			weekEnd := weekStart.AddDate(0, 0, 6)
			if weekEnd.After(end) {
				weekEnd = end
			}
			periods = append(periods, periodRange{
				Label: weekStart.Format("2006-W01"),
				Start: weekStart,
				End:   weekEnd,
			})
			current = weekEnd.AddDate(0, 0, 1)
		}
	case Monthly:
		for current.Before(end) || current.Equal(end) {
			monthStart := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())
			monthEnd := monthStart.AddDate(0, 1, -1)
			if monthEnd.After(end) {
				monthEnd = end
			}
			periods = append(periods, periodRange{
				Label: monthStart.Format("2006-01"),
				Start: monthStart,
				End:   monthEnd,
			})
			current = monthEnd.AddDate(0, 0, 1)
		}
	case Yearly:
		for current.Before(end) || current.Equal(end) {
			yearStart := time.Date(current.Year(), 1, 1, 0, 0, 0, 0, current.Location())
			yearEnd := time.Date(current.Year(), 12, 31, 0, 0, 0, 0, current.Location())
			if yearEnd.After(end) {
				yearEnd = end
			}
			periods = append(periods, periodRange{
				Label: yearStart.Format("2006"),
				Start: yearStart,
				End:   yearEnd,
			})
			current = yearEnd.AddDate(0, 0, 1)
		}
	default: // Daily
		for current.Before(end) || current.Equal(end) {
			periods = append(periods, periodRange{
				Label: current.Format("2006-01-02"),
				Start: current,
				End:   current,
			})
			current = current.AddDate(0, 0, 1)
		}
	}

	return periods
}

func (s *StatisticService) getWeekStart(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	return t.AddDate(0, 0, -(weekday - 1))
}

func (s *StatisticService) buildStatisticCondition(userId uint, startTime, endTime time.Time, categoryIds []uint) *transactionModel.StatisticCondition {
	condition := transactionModel.NewStatisticConditionBuilder(userId).
		WithDate(startTime, endTime)

	if len(categoryIds) > 0 {
		condition = condition.WithCategoryIds(categoryIds)
	}

	return condition.Build()
}
