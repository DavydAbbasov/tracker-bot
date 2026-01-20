package subscription

// ---------------------------------------------------------------------
// Inline callbacks (actions)
const (
	SubscriptionCBTariffPlans   = "subscription:tariff:plans"
	SubscriptionCBFreePlan      = "subscription:free:plan"
	SubscriptionCBSupport       = "subscription:support"
	SubscriptionCBPaymentChange = "subscription:payment:change"
)

// ---------------------------------------------------------------------
// Buttons (Inline + Reply)

// Entry inline menu buttons
const (
	SubscriptionButtonTariffPlans   = "🗓 Tariff plans"
	SubscriptionButtonFreePlan      = "🎁 Free"
	SubscriptionButtonSupport       = "🛫 Support"
	SubscriptionButtonPaymentChange = "💳 Change payment"
)

// ---------------------------------------------------------------------
// Track UI texts (titles/labels shown inside messages)

// Main screen
const (
	SubscriptionUIMainTitle      = "💳 Subscription"
	SubscriptionUIMainTariffPlan = "🗓 Tariff plan:"
	SubscriptionUIMainDaysEnd    = "🕐 Days end:"
	SubscriptionUIMainMessage    = "To subscribe, go to: 🗓 Tariff plans"
)
