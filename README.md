# Tiimit - Distribute football teams

Tiimit is tool to distribute football teams easily for the team leader.

## Install and deploy

Deploying this app to Raspberry Pi can be done by:

Build the production binaries
```
make prod/build-arm64
```

Copy the binaries to the server
```
scp /build/web username@server:/opt/tiimit/
```

Set the environment variables in `.env` in the project root

`DB_PATH=path-to-db`