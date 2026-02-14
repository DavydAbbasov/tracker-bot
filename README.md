tracker-bot is a Telegram bot for tracking your daily activities over configurable time intervals.
It lets you create activities, receive periodic prompts, select what you are currently doing, and automatically record time spent on each activity.
Over time, the bot generates detailed reports showing how much time you spent on specific activities per day or aggregated across multiple days.

## Recent Improvements

The project was recently cleaned up and aligned to a clearer architecture without changing core behavior.

- `application` was split into clear phases:
  - config is parsed in `cmd/tracker-bot/main.go`
  - dependencies are wired in `Application.Build(...)`
  - runtime is started in `Application.Run(...)`
- handler module naming was clarified:
  - `internal/handlers/router.go` was renamed to `internal/handlers/module.go`
- dispatcher and handlers were refactored for readability:
  - less duplicate code
  - clearer helper methods
  - short, practical function comments
- model layer was cleaned:
  - removed unused constants
  - clarified DTO comments
- utility layer was normalized:
  - `internal/utils/tgcient` renamed to `internal/utils/tgclient`
  - dead code removed from Telegram client helper
  - PostgreSQL client naming simplified
- Docker setup improved:
  - removed `.env` copy into image
  - non-root runtime user
  - better multi-arch build support
- developer experience improved:
  - simplified `Makefile` with readable targets
  - `.DS_Store` is now ignored in Git
