package learning

// ---------------------------------------------------------------------
// Inline callbacks (actions)

const (
	LearningCBAddCollection    = "learning:add:collection"
	LearningCBRandomWords      = "learning:random:words"
	LearningCBSwitchCollection = "learning:switch:collection"
	LearningCBSummaryLearning  = "learning:summary:learning"
	LearningCBBaseWords        = "learning:base:words"
)

// ---------------------------------------------------------------------
// Buttons (Inline + Reply)

// Entry inline menu buttons
const (
	LearningButtonAddCollection    = "➕ Создать подборку"
	LearningButtonRandomWords      = "🎲 Случайная подборка"
	LearningButtonSwitchCollection = "🔁 Архив подборок"
	LearningButtonSummaryLearning  = "📈 Статистика"
	LearningButtonBaseWords        = "🗂 База слов"
)

// Add Collection reply menu buttons
const (
	LearningButtonHelp = "ℹ️ Помощь"
	LearningButtonHome = "🏠 Home"
)

// Add words reply menu buttons
const (
	LearningButtonAddWord  = "➕ Добавить слово"
	LearningButtonComplete = "✅ Завершить"
	LearningButtonBackHome = "🏠 Home"
)
