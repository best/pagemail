#!/bin/sh
set -eu

APP_USER="pagemail"
PUID="${PUID:-1000}"
PGID="${PGID:-1000}"

# Validate PUID/PGID are numeric and not root (0)
validate_id() {
  case "$1" in
    ''|*[!0-9]*) echo "ERROR: $2 must be numeric, got '$1'" >&2; exit 1 ;;
  esac
  # Use arithmetic comparison to handle "0", "00", "000" etc.
  if [ "$1" -eq 0 ]; then
    echo "ERROR: $2=0 (root) is not allowed for security reasons" >&2
    exit 1
  fi
}

validate_id "$PUID" "PUID"
validate_id "$PGID" "PGID"

if [ "$(id -u)" = "0" ]; then
  current_uid="$(id -u "$APP_USER")"
  current_gid="$(id -g "$APP_USER")"

  # Adjust GID if different
  if [ "$PGID" != "$current_gid" ]; then
    echo "Updating $APP_USER GID: $current_gid -> $PGID"
    groupmod -o -g "$PGID" "$APP_USER"
  fi

  # Adjust UID if different
  if [ "$PUID" != "$current_uid" ]; then
    echo "Updating $APP_USER UID: $current_uid -> $PUID"
    usermod -o -u "$PUID" "$APP_USER"
  fi

  # Fix ownership of app directory and internal paths (-h: don't follow symlinks)
  if ! chown -Rh "$APP_USER:$APP_USER" /app /var/log/nginx /var/run/nginx /var/log/supervisor /tmp/crashpad 2>&1; then
    echo "WARN: Some paths could not be chowned; pagemail may have permission issues" >&2
  fi
else
  echo "WARN: Not running as root, skipping PUID/PGID adjustments" >&2
fi

exec "$@"
