# Tiimit - Distribute football teams

Tiimit is a webapp to distribute the players into even teams in football practices. It helps the coach of the team to set up fun and even games for everyone with minimal effort.

## How this was built

Tiimit started as my programming learning project, and up to November 2026 I wrote the code by hand. I used Claude as a mentor, asking it to guide me like a senior engineer would, but it didn't write the code. That version is preserved at the [`handwritten-final`](https://github.com/Zigelzi/go-tiimit/tree/handwritten-final) tag.

From November 2026, Tiimit is built with an agentic workflow, with Claude Code writing most of the code. I own the product decisions, the specs, the architecture rules and the acceptance of every change. Each change starts from a written spec with acceptance conditions, gets tests written before the code, and is reviewed before it merges. The workflow is described in [`spec/constitution.md`](spec/constitution.md) and the code's structure in [`ARCHITECTURE.md`](ARCHITECTURE.md). Commits made with the agent carry a `Co-Authored-By: Claude` line.

## Features

### Manage players

1. Add a player by hand, for example a mid-season joiner.
2. Adjust the scores and goalie willingness of players to ensure balanced teams.
3. Set a player inactive when they leave or pause.

### Distribute teams

1. Distribute players to the teams for a practice from attendace report exported from MyClub.
2. View existing practices.
3. Mark players who have vest to distribute the vests easily in the beginning of the practice.
4. Move players in a practice to rebalance the teams.

### User management

1. Invite new organizers with a one-time invitation; they choose their own username and password. Every page with player data requires login.
2. On an empty database (a fresh install or a recovery), the first user is created on startup from the `FIRST_USER_USERNAME` and `FIRST_USER_PASSWORD` settings. See `.env.example`.

## Installation and deployment

The webapp is deployed to a Raspberry Pi on the same network as the development computer with `make deploy`. It runs the tests, builds the app, backs up the database on the Pi, installs the new version and checks that it answers.

See `docs/deployment-guide.md` for the one-time setup and `deploy/deploy.conf.example` for the settings.
