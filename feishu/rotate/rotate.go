package rotate

import (
	"fmt"
	"github.com/prometheus/common/model"
	"regexp"
	"strings"
	"time"
)

type MentionRotator struct {
	baseDate  time.Time
	cycleDays int
	openIDs   []string
}

// by now, support week and day
var durationRE = regexp.MustCompile(`^(([0-9]+)w)?(([0-9]+)d)?$`)

// parse to days
func parseDuration(durationStr string) (int, error) {
	durationStr = strings.TrimSpace(durationStr)
	matches := durationRE.FindStringSubmatch(durationStr)
	if matches == nil {
		return 0, fmt.Errorf("not a valid duration string: %q", durationStr)
	}

	duration, err := model.ParseDuration(durationStr)
	if err != nil {
		return 0, err
	}
	return int(time.Duration(duration) / time.Millisecond / (1000 * 60 * 60 * 24)), nil
}

func New(rotationStr string, openIDs []string) (*MentionRotator, error) {
	parts := strings.Split(rotationStr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid rotation string: %v", rotationStr)
	}
	baseDateStr := parts[0]
	rotateDuration := parts[1]

	// parse base date: add timezone (system timezone)
	baseDateStr += time.Now().Format("Z07:00")
	baseDate, err := time.Parse("2006-01-02Z07:00", baseDateStr)
	if err != nil {
		return nil, err
	}

	// parse rotate duration
	days, err := parseDuration(rotateDuration)
	if err != nil {
		return nil, err
	}
	if days <= 0 {
		return nil, fmt.Errorf("rotate duration at least 1: %v", days)
	}

	return &MentionRotator{
		baseDate:  baseDate,
		cycleDays: days,
		openIDs:   openIDs,
	}, nil
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

func adjustDays(relativeDays int, cycleDays int) int {
	if cycleDays <= 0 {
		panic("unexpected")
	}
	if relativeDays < 0 {
		// example: -3 -2 -1 1 2 3 4 5
		// cycle = 2
		// -1 => 3
		relativeDays = abs(cycleDays - relativeDays)
	} else {
		// = 0 means at that day
		relativeDays += 1
	}
	return relativeDays
}

func (r MentionRotator) Rotate(t time.Time) []string {
	if len(r.openIDs) <= 1 {
		return r.openIDs
	}
	days := int(t.Sub(r.baseDate) / time.Hour / 24)
	days = adjustDays(days, r.cycleDays)
	index := (bucketIndexEveryN(days, r.cycleDays) - 1) % len(r.openIDs)
	var res []string
	res = append(res, r.openIDs[index])
	return res
}

func bucketIndexEveryN(v, bucketSize int) int {
	if bucketSize <= 0 {
		panic("unexpected")
	}
	if v <= 0 {
		panic("unexpected")
	}
	if v%bucketSize != 0 {
		return (v - v%bucketSize + bucketSize) / bucketSize
	}
	return v / bucketSize
}
