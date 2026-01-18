package track

// ---------------------------------------------------------------------
// Inline callbacks (actions)
const (
	TrackCBActivitySelect     = "track:activity:select"
	TrackCBActivityCreate     = "track:activity:create"
	TrackCBReportSummary      = "track:report:summary"
	TrackCBArchiveOpen        = "track:archive:open"
	TrackCBActivityReportOpen = "track:activity:report"
)

// ---------------------------------------------------------------------
// Buttons (Inline + Reply)

// Entry inline menu buttons
const (
	TrackButtonSelectActivity = "📂 Activities"
	TrackButtonCreateActivity = "➕ New Activity"
	TrackButtonExitTracking   = "⏹ Stop Tracking"
	TrackButtonViewReports    = "📈 Reports"
	TrackButtonViewArchive    = "🗄 Archive"
)

// Common reply buttons
const (
	TrackButtonToday    = "📊 Today"
	TrackButtonBack     = "◀ Back"
	TrackButtonBackHome = "🏠 Home"
)

// Report reply menu buttons
const (
	TrackButtonReportPeriod = "📅 Period"
	TrackButtonReportWeek   = "🗓 Week"
	TrackButtonReportExport = "📤 Export"
	TrackButtonReportDelete = "🗑 Delete"
)

// Activity manage reply menu buttons
const (
	TrackButtonActivityActivate = "📳 Activate"
	TrackButtonActivityArchive  = "🛒 Archive"
	TrackButtonActivityDelete   = "🗑 Delete"
)

// Timer reply menu buttons
const (
	TrackButtonTimer15     = "⏱ 15 min"
	TrackButtonTimer60     = "⏱ 60 min"
	TrackButtonTimerCreate = "➕ Custom Timer"
)

// ---------------------------------------------------------------------
// Track UI texts (titles/labels shown inside messages)

// Main screen
const (
	TrackUIMainTitle                = "📈 Tracking"
	TrackUIMainLabelCurrentActivity = "📌 Current activity:"
	TrackUIMainLabelTodayTime       = "⏱ Tracked today:"
	TrackUIMainLabelStreak          = "🔥 Streak:"
	TrackUIMainLabelTodayCount      = "✅ Sessions today:"
)

// Activity report screen
const (
	TrackUIReportTitle                = "📌 Activity report"
	TrackUIReportLabelStartDate       = "📅 Started:"
	TrackUIReportLabelConsecutiveDays = "📈 Streak:"
	TrackUIReportLabelTodayTimeTotal  = "⏱ Today total:"
	TrackUIReportLabelAvgDailyTime    = "📊 Daily average:"
	TrackUIReportLabelTodayDate       = "🗓 Date:"
)

// ---------------------------------------------------------------------
// Messages (plain texts, not labels/titles)
const (
	TrackMsgActivityListTitle     = "📂 Select Activity"
	TrackMsgActivityListConfirmed = "📂 Activated Activities:"
)
