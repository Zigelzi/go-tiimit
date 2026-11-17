# Tiimit - Distribute football teams

Tiimit is a webapp to distribute the players into even teams in football practices. It helps the coach of the team to set up fun and even games for everyone with minimal effort. Tiimit comes currently with a complementary CLI that is primarily used to manage the users of the app internally.

## How this was built

Tiimit started as my programming learning project, and up to November 2026 I wrote the code by hand. I used Claude as a mentor, asking it to guide me like a senior engineer would, but it didn't write the code. That version is preserved at the [`handwritten-final`](https://github.com/Zigelzi/go-tiimit/tree/handwritten-final) tag.

From November 2026, Tiimit is built with an agentic workflow, with Claude Code writing most of the code. I own the product decisions, the specs, the architecture rules and the acceptance of every change. Each change starts from a written spec with acceptance conditions, gets tests written before the code, and is reviewed before it merges. The workflow is described in [`spec/constitution.md`](spec/constitution.md) and the code's structure in [`ARCHITECTURE.md`](ARCHITECTURE.md). Commits made with the agent carry a `Co-Authored-By: Claude` line.

## Features

### Manage players

1. Import players from MyClub to Tiimit by using Excel exported from MyClub.
2. Add a player by hand, for example a mid-season joiner.
3. Adjust the scores and goalie willingness of players to ensure balanced teams.
4. Set a player inactive when they leave or pause.

### Distribute teams

1. Distribute players to the teams for a practice from attendace report exported from MyClub.
2. View existing practices.
3. Mark players who have vest to distribute the vests easily in the beginning of the practice.
4. Move players in a practice to rebalance the teams.

### User management (CLI)

1. Register new users to the app to prevent unauthorized access to the app.

## Installation and deployment

The webapp can be currently deployed to Raspberry Pi 3+ by using the `deploy.sh` command. It will build and deploy the webapp to Raspberry Pi connected to the same network as the development computer.

See `deploy.sh` and `.deploy.conf.example` for details.
