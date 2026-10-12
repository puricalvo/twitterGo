#!/bin/sh
set -eu

APP_DIR="${APP_DIR:-/home/Cadesa/twittergo-api}"
PIDFILE="$APP_DIR/api.pid"
LOGFILE="$APP_DIR/api.log"

cd "$APP_DIR"
set -a
. "$APP_DIR/api.env"
set +a

: "${MONGODB_URI:?Falta MONGODB_URI}"
: "${MONGODB_DATABASE:?Falta MONGODB_DATABASE}"
: "${JWT_SIGN:?Falta JWT_SIGN}"
: "${PORT:?Falta PORT}"
: "${CORS_ORIGIN:?Falta CORS_ORIGIN}"

if [ -f "$PIDFILE" ]; then
    pid=$(cat "$PIDFILE")
    case "$pid" in
        ''|*[!0-9]*) pid="" ;;
    esac
    if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
        echo "La API ya está ejecutándose (PID $pid)."
        exit 0
    fi
    rm -f "$PIDFILE"
fi

nohup "$APP_DIR/twittergo-freebsd-amd64" >> "$LOGFILE" 2>&1 </dev/null &
pid=$!
echo "$pid" > "$PIDFILE"
sleep 2

if kill -0 "$pid" 2>/dev/null; then
    echo "API iniciada (PID $pid). Registro: $LOGFILE"
else
    rm -f "$PIDFILE"
    echo "La API no pudo arrancar. Últimas líneas del registro:"
    tail -n 30 "$LOGFILE"
    exit 1
fi
