package learning

import (
	"fmt"
	"tracker-bot/internal/models"
)

func LearningMenuText(stats models.LearningStats) string {
	return fmt.Sprintf(
		"%s\n\n%s *%s*\n%s *%d*\n%s *%d*\n%s *%d*\n%s *%s*\n",
		LearningUIMainTitle,
		LearningUIMainLanguage, stats.Language,
		LearningUIMainTotalWords, stats.TotalWords,
		LearningUIMainTodayWords, stats.TodayWords,
		LearningUIMainLearnedWords, stats.LearnedWords,
		LearningUIMainNextWordIn, stats.NextWordIn,
	)
}
