package habitcompletion

import (
	kaievolution "kai-back/internal/modules/kai/evolution"

	"strings"
	"unicode"
)

const (
	kaiStateCurious = "curioso"
	kaiStateHappy   = "feliz"
	kaiStateProud   = "orgulloso"

	kaiModeSoft        = "suave"
	kaiModeReflective  = "reflexivo"
	kaiModeEnergetic   = "energico"
	kaiModeCelebration = "celebracion"

	kaiEnergyMin = 10
	kaiEnergyMax = 100
	kaiBondMin   = 0
	kaiBondMax   = 100

	happyCompletedHabitsThreshold = 3
	proudStreakThreshold          = 7
)

type kaiCompletionRules struct {
	currentState string
	currentStage string
	currentMode  string
	energyGain   int
	bondGain     int
	stageChanged bool
}

func calculateStage(totalXP int) string {
	return kaievolution.StageForXP(totalXP)
}

func calculateState(completedToday int, currentStreak int) string {
	if currentStreak >= proudStreakThreshold {
		return kaiStateProud
	}

	if completedToday >= happyCompletedHabitsThreshold {
		return kaiStateHappy
	}

	return kaiStateCurious
}

func calculateMode(attributeName string, currentStreak int, stageChanged bool) string {
	if currentStreak >= proudStreakThreshold || stageChanged {
		return kaiModeCelebration
	}

	normalized := normalizeRuleText(attributeName)

	switch normalized {
	case "sabiduria", "conciencia":
		return kaiModeReflective
	case "vitalidad", "fuerza":
		return kaiModeEnergetic
	default:
		return kaiModeSoft
	}
}

func calculateEnergyGain(difficulty *string, currentStreak int) int {
	gain := 3

	if difficulty != nil {
		switch normalizeRuleText(*difficulty) {
		case "alta", "alto", "dificil", "dificultad alta":
			gain += 5
		}
	}

	if currentStreak > 1 {
		gain += 2
	}

	return gain
}

func calculateBondGain(currentStreak int) int {
	gain := 1

	if currentStreak >= proudStreakThreshold {
		gain += 2
	}

	return gain
}

func clampInt(value int, min int, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

func normalizeRuleText(value string) string {
	lower := strings.ToLower(strings.TrimSpace(value))

	replacer := strings.NewReplacer(
		"á", "a",
		"é", "e",
		"í", "i",
		"ó", "o",
		"ú", "u",
		"ü", "u",
		"ñ", "n",
	)

	normalized := replacer.Replace(lower)

	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, normalized)
}

func messageFiltersForRules(state string, mode string) ([]string, []string) {
	switch {
	case state == kaiStateProud || mode == kaiModeCelebration:
		return []string{"racha", "motivacion", "motivación"}, []string{"racha", "celebracion", "celebración"}
	case state == kaiStateHappy || mode == kaiModeEnergetic:
		return []string{"motivacion", "motivación", "autocuidado"}, []string{"motivacion", "motivación", "habito", "hábito"}
	case mode == kaiModeReflective:
		return []string{"reflexion", "reflexión", "estoico", "conciencia"}, []string{"reflexion", "reflexión", "conciencia"}
	default:
		return []string{"motivacion", "motivación", "autocuidado"}, []string{"general", "habito", "hábito"}
	}
}
