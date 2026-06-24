package evolution

import (
	"strings"
	"time"
)

const (
	StagePuppy = "cachorro"
	StageYoung = "joven"
	StageAdult = "adulto"

	YoungMinXP = 30
	AdultMinXP = 60

	EventDuration = 24 * time.Hour
)

func StageForXP(totalXP int) string {
	if totalXP >= AdultMinXP {
		return StageAdult
	}
	if totalXP >= YoungMinXP {
		return StageYoung
	}
	return StagePuppy
}

func NormalizeStage(stage string) string {
	switch strings.ToLower(strings.TrimSpace(stage)) {
	case "bebe", StagePuppy:
		return StagePuppy
	case StageYoung:
		return StageYoung
	case StageAdult:
		return StageAdult
	default:
		return strings.ToLower(strings.TrimSpace(stage))
	}
}

func ImageKeyForStage(stage string) string {
	switch NormalizeStage(stage) {
	case StageYoung:
		return "kai8"
	case StageAdult:
		return "kai12"
	default:
		return "kai1"
	}
}

func EventWindow(lastEvolution *time.Time, now time.Time) (bool, *time.Time) {
	if lastEvolution == nil {
		return false, nil
	}

	expiresAt := lastEvolution.Add(EventDuration)
	active := !now.Before(*lastEvolution) && now.Before(expiresAt)
	return active, &expiresAt
}
