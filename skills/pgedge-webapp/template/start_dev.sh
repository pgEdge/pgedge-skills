#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$SCRIPT_DIR/bin"
CLIENT_DIR="$SCRIPT_DIR/client"
SERVER_SRC_DIR="$SCRIPT_DIR/server/src"
SERVER_CONFIG="$BIN_DIR/<BINARY_NAME>.yaml"
EXAMPLE_CONFIG="$SCRIPT_DIR/examples/<BINARY_NAME>.yaml"
SERVER_BIN="$BIN_DIR/<BINARY_NAME>"
SERVER_LOG="/tmp/<BINARY_NAME>.log"
VITE_LOG="/tmp/<BINARY_NAME>-vite.log"

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; RED='\033[0;31m'; NC='\033[0m'

cleanup() {
    [ -n "${CLIENT_PID:-}" ] && kill "$CLIENT_PID" 2>/dev/null || true
    [ -n "${SERVER_PID:-}" ] && kill "$SERVER_PID" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

mkdir -p "$BIN_DIR"

if [ ! -f "$SERVER_CONFIG" ]; then
    echo -e "${YELLOW}No config found, copying example to $SERVER_CONFIG${NC}"
    cp "$EXAMPLE_CONFIG" "$SERVER_CONFIG"
fi

if [ ! -f "$SERVER_BIN" ] || \
   [ -n "$(find "$SERVER_SRC_DIR" -name '*.go' -newer "$SERVER_BIN" 2>/dev/null)" ]; then
    echo -e "${BLUE}Building server...${NC}"
    (cd "$SERVER_SRC_DIR" && go build -o "$SERVER_BIN" "./cmd/<BINARY_NAME>")
fi

if [ ! -d "$CLIENT_DIR/node_modules" ] || \
   [ "$CLIENT_DIR/package.json" -nt "$CLIENT_DIR/node_modules" ]; then
    echo -e "${BLUE}Installing client dependencies...${NC}"
    (cd "$CLIENT_DIR" && npm install)
fi

if [ ! -f "$BIN_DIR/auth.db" ]; then
    echo -e "${YELLOW}Initializing auth database (seeding admin)...${NC}"
    "$SERVER_BIN" --init --config "$SERVER_CONFIG"
fi

echo -e "${GREEN}Starting server...${NC}"
"$SERVER_BIN" --config "$SERVER_CONFIG" > "$SERVER_LOG" 2>&1 &
SERVER_PID=$!

for _ in $(seq 1 30); do
    if curl -fs "http://localhost:<HTTP_PORT>/health" >/dev/null 2>&1; then break; fi
    if ! kill -0 "$SERVER_PID" 2>/dev/null; then
        echo -e "${RED}Server died during startup${NC}"
        tail -20 "$SERVER_LOG"
        exit 1
    fi
    sleep 1
done

echo -e "${GREEN}Starting Vite...${NC}"
(cd "$CLIENT_DIR" && npm run dev > "$VITE_LOG" 2>&1) &
CLIENT_PID=$!

for _ in $(seq 1 30); do
    if grep -q "Local:" "$VITE_LOG" 2>/dev/null; then break; fi
    sleep 1
done

cat <<EOF
${GREEN}<PROJECT_NAME> dev environment running.${NC}
  Server:   http://localhost:<HTTP_PORT>  (log: $SERVER_LOG)
  Frontend: http://localhost:5173         (log: $VITE_LOG)
  Admin:    admin (rotate the seed password after first login)

Press Ctrl+C to stop.
EOF

wait "$CLIENT_PID"
