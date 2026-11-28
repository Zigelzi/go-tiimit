#!/usr/bin/env bash
# deploy.sh tests Tiumo, builds it for the Pi, copies it there and has the
# Pi's root script back up the database, install it and check that it
# answers. Run it with `make deploy`. docs/deployment-guide.md has the
# one-time setup it relies on.
#
# Only committed work on main is deployed, so the Pi's log can name the commit
# that is running. ALLOW_BRANCH=1 lets a committed branch through, for trying
# a milestone on the phones before it merges, but only when it changes no
# migration file compared with main: a migration stays in the live database
# even if the branch is changed or dropped later.
set -euo pipefail

cd "$(dirname "$0")/.."

conf=deploy/deploy.conf
if [ ! -f "$conf" ]; then
	echo "$conf not found: copy it from $conf.example and fill it in" >&2
	exit 1
fi
# shellcheck source=deploy.conf.example
source "$conf"

pi="$PI_USER@$PI_HOST"
binary=build/web

step=""
start() {
	step=$1
	echo
	echo "==> $step"
}
trap 'echo; echo "DEPLOY FAILED while $step" >&2' ERR

start "checking the work is committed"
if [ -n "$(git status --porcelain)" ]; then
	echo "There are uncommitted changes or new files; commit or remove them first:" >&2
	git status --short >&2
	false
fi
branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$branch" != main ]; then
	if [ "${ALLOW_BRANCH:-}" != 1 ]; then
		echo "This is $branch, not main. Merge it first, or run ALLOW_BRANCH=1 make deploy to try it on the phones." >&2
		false
	fi
	if ! git diff --quiet main HEAD -- sql/schema; then
		echo "$branch changes migration files compared with main:" >&2
		git diff --stat main HEAD -- sql/schema >&2
		echo "A migration stays in the live database even if the branch changes later, so merge it to main first." >&2
		false
	fi
	echo "WARNING: deploying $branch, not main"
fi
commit=$(git rev-parse HEAD)
echo "Deploying $commit ($branch)"

start "running the tests"
make test

start "checking the Pi is 64-bit ARM"
arch=$(ssh "$pi" uname -m)
if [ "$arch" != aarch64 ]; then
	echo "The Pi reports $arch, but the build is for aarch64" >&2
	false
fi

start "building for the Pi"
make prod/build-arm64

# Go records in the binary whether the tree had changes. This also catches
# make generate rewriting committed files, which the first check cannot see.
start "checking the build matches the commit"
if ! go version -m "$binary" | grep -q 'vcs.modified=false'; then
	echo "The build includes changes that are not committed:" >&2
	git status --short >&2
	false
fi

start "copying the build to the Pi"
scp -q "$binary" "$pi:tiimit-new"

start "backing up, installing and checking on the Pi"
status=0
ssh "$pi" sudo -n /usr/local/sbin/tiimit-pi install || status=$?

case $status in
0) ;;
3)
	cat >&2 <<'EOF'

DEPLOY FAILED: the new version is installed but does not answer.

See why, on the Pi:
  sudo journalctl -u tiimit -n 50 --no-pager

Put the previous version back, on the Pi:
  sudo systemctl stop tiimit
  sudo install -o root -g root -m 755 /opt/tiimit/web.previous /opt/tiimit/web

If this version added a database migration, also put back the deploy backup
from just before the first failed deploy, chosen by the time in its name,
before starting the app: "Restoring a backup" in docs/deployment-guide.md.
Anything saved since that backup is lost. Otherwise start the app:
  sudo systemctl reset-failed tiimit
  sudo systemctl start tiimit
EOF
	exit 1
	;;
255)
	cat >&2 <<'EOF'

DEPLOY FAILED: the connection to the Pi was lost.

If it dropped after the copy, the Pi finishes the deploy on its own. See how
it ended, on the Pi:
  sudo journalctl -t tiimit-pi -n 20 --no-pager

If that ends with "FAILED", follow "Rollback" in docs/deployment-guide.md.
Otherwise run make deploy again.
EOF
	exit 1
	;;
*)
	echo >&2
	echo "DEPLOY FAILED: the Pi stopped the deploy; its message above says why." >&2
	exit 1
	;;
esac

# The check that counts ran on the Pi. This one also goes through DNS, TLS
# and nginx, so it can fail while the app is fine, and only warns.
start "checking $PUBLIC_URL from here"
if code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$PUBLIC_URL") && [ "$code" = 200 ]; then
	echo "$PUBLIC_URL answers"
else
	echo "WARNING: $PUBLIC_URL answered ${code:-nothing}, not a normal page. The app works on the Pi, so look at nginx, DNS or the network."
fi

echo
echo "Deployed $commit"
