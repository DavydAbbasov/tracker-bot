package learning

import (
	"fmt"
	"tracker-bot/internal/models"
)

func LearningMenuText(stats models.LearningStats) string {
	return fmt.Sprintf(
		"%s\n\n%s *%s*\n%d *%d*\n%d *\n",
		LearningUIMainTitle,
		LearningUIMainLanguage, stats.Language,
		LearningUIMainTotalWords, stats.TotalWords,
		LearningUIMainTodayWords, stats.TodayWords,
		LearningUIMainLearnedWords, stats.LearnedWords,
		LearningUIMainNextWordIn, stats.NextWordIn,
	)
}
