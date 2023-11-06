package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

func ApplyFilterCondition(query *gorm.DB, filter map[string]interface{}) (*gorm.DB, error) {
	for condition, value := range filter {
		if value == nil {
			query = query.Where(fmt.Sprintf("%s IS NULL", condition))
			continue
		}

		// if value can be casted to int, then apply int filter condition
		isInt, err := strconv.Atoi(fmt.Sprintf("%v", value))
		if err == nil {
			query = query.Where(fmt.Sprintf("%s = ?", condition), isInt)
			continue
		}

		switch v := value.(type) {
		case string:
			if dateRange, err := parseDateRange(v); err == nil {
				query = query.Where(fmt.Sprintf("%s >= ? AND %s <= ?", condition, condition), dateRange[0], dateRange[1])
			} else if isValidOperatorCondition(condition) {
				query = query.Where(fmt.Sprintf("%s %s ?", condition, v))
			} else {
				query = applyStringFilterCondition(query, condition, v)
			}
		case int:
			query = query.Where(fmt.Sprintf("%s = ?", condition), v)
		default:
			return nil, fmt.Errorf("unsupported filter type for %s: %T", condition, value)
		}
	}

	return query, nil
}

func parseDateRange(value string) ([]string, error) {
	dateRange := strings.Split(value, "~")
	if len(dateRange) != 2 {
		return nil, fmt.Errorf("invalid date range format")
	}

	start := dateRange[0]
	end := dateRange[1]

	if isValidDate(start) && isValidDate(end) {
		// check if start date and end date is date only or date with time
		if !strings.Contains(start, " ") {
			start += " 00:00:00"
		}
		if !strings.Contains(end, " ") {
			end += " 23:59:59"
		}
	}

	return []string{start, end}, nil
}

func isValidOperatorCondition(condition string) bool {
	supportedOperators := map[string]bool{
		">=": true,
		">":  true,
		"<=": true,
		"<":  true,
	}
	parts := strings.SplitN(condition, " ", 2)
	if len(parts) != 2 {
		return false
	}
	operator := parts[1]
	return supportedOperators[operator]
}

func applyStringFilterCondition(query *gorm.DB, condition string, value string) *gorm.DB {
	if isValidDate(value) {
		query = query.Where(fmt.Sprintf("DATE(%s) = ?", condition), value)
	} else {
		query = query.Where(fmt.Sprintf("%s::text ILIKE ?", condition), fmt.Sprintf("%%%s%%", value))
	}
	return query
}

func isValidDate(value string) bool {
	dateFormats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"02/01/2006",
		"02/01/2006 15:04:05",
		"01-02-2006",
		"01-02-2006 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		"01/02/2006",
		"01/02/2006 15:04:05",
		"02-Jan-2006",
		"02-Jan-2006 15:04:05",
		"Jan 02, 2006",
		"Jan 02, 2006 15:04:05",
		"January 02, 2006",
		"January 02, 2006 15:04:05",
		"02-January-2006",
		"02-January-2006 15:04:05",
	}

	for _, format := range dateFormats {
		_, err := time.Parse(format, value)
		if err == nil {
			return true
		}
	}
	return false
}
