tracker-bot is a Telegram bot for tracking your daily activities over configurable time intervals.
It lets you create activities, receive periodic prompts, select what you are currently doing, and automatically record time spent on each activity.
Over time, the bot generates detailed reports showing how much time you spent on specific activities per day or aggregated across multiple days.

## What You Can Do

- Create and manage activities (`active` / `archived`)
- Select activities you want to track right now
- Start a timer with fixed interval prompts
- Answer prompt messages and automatically save tracked time
- Get statistics for:
  - today
  - custom date periods
  - selected activities only
- View reports in:
  - text format
  - chart-like format

## How Tracking Works

1. Create activities (for example: Go, English, Workout).
2. Select active activities.
3. Start timer prompts (e.g. every 15/30 minutes).
4. When bot asks "What are you doing now?", choose one activity.
5. Bot saves a retro session for the last interval and builds reports from saved sessions.
