# Tiimit - Distribute football teams

Tiimit is webapp to distribute football teams easily for the team leader. It comes with complementary CLI that is primarily used to manage the users of the app internally. The CLI is also used for some club management activities.

## Features

### Manage players

1. Import players from MyClub to Tiimit by using Excel exported from MyClub.
2. Manage the scores of players used to distribute the teams.

### Distribute teams

1. Distribute players to the teams for a practice from attendace report exported from MyClub.
2. View existing practices.
3. Mark players who have vest to distribute the vests easily in the beginning of the practice.
4. Move players in a practice to rebalance the teams.

### User management (CLI)

1. Register new users to the app to prevent unauthorized access to the app.

## Installation and deployment

The webapp can be currently deployed to Raspberry Pi 3+ by using the `deploy.sh` command. It will build and deploy the webapp to Raspberry Pi connected to the same network as the development computer.

See `deploy.sh` and `deploy.conf.example` for details.
