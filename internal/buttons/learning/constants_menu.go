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
	LearningButtonAddCollection    = "➕ Create a collection"
	LearningButtonRandomWords      = "🎲 Random collection"
	LearningButtonSwitchCollection = "🔁 Archive of collections"
	LearningButtonSummaryLearning  = "📈 Statistics"
	LearningButtonBaseWords        = "🗂 Word base"
)

// Add Collection reply menu buttons
const (
	LearningButtonHelp = "ℹ️ Help"
	LearningButtonHome = "🏠 Home"
)

// Add words reply menu buttons
const (
	LearningButtonAddWord  = "➕ Add a word"
	LearningButtonComplete = "✅ Finish"
	LearningButtonBackHome = "🏠 Home"
)

// ---------------------------------------------------------------------
// Track UI texts (titles/labels shown inside messages)

// Main screen
const (
	LearningUIMainTitle        = "🧠 Learning"
	LearningUIMainLanguage     = "🌐 Language:"
	LearningUIMainTotalWords   = "📊 Total Words:"
	LearningUIMainTodayWords   = "📘 Today Words:"
	LearningUIMainLearnedWords = "✅ Learned Words:"
	LearningUIMainNextWordIn   = "🕐 Next Word In:"
)
