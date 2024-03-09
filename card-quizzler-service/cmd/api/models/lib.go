package models

const (
	maxDifficultThreshold = 5 // Maximum negative answer streak before marking a term as difficult
	maxRecoveryThreshold  = 5 // Maximum positive answer streak before unmarking a term as difficult
)

// UpdateStreaksAndUpdateDifficulty updates the streaks and difficulty status of a term based on the provided answer.
func UpdateStreaksAndUpdateDifficulty(term *Term, answer bool) {
	// Update streaks based on the answer
	switch {
	case answer:
		term.NegativeAnswerStreak = 0
		term.PositiveAnswerStreak = term.PositiveAnswerStreak + 1
	default:
		term.NegativeAnswerStreak = term.NegativeAnswerStreak + 1
		term.PositiveAnswerStreak = 0
	}

	// Update difficulty status
	switch term.IsDifficult {
	case true:
		if term.PositiveAnswerStreak >= maxRecoveryThreshold {
			term.IsDifficult = false
			term.PositiveAnswerStreak = 0
			term.NegativeAnswerStreak = 0
		}
	default:
		if term.NegativeAnswerStreak >= maxDifficultThreshold {
			term.IsDifficult = true
			term.NegativeAnswerStreak = 0
			term.PositiveAnswerStreak = 0
		}
	}
}
