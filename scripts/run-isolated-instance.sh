#!/bin/sh
# Start an isolated PicoClaw instance.
#
# Each instance gets its own PICOCLAW_HOME, config, workspace, logs, pid file,
# skills, launcher auth store, and default cache/state directories.
#
# Usage:
#   ./scripts/run-isolated-instance.sh <name> [gateway|launcher] [-- extra args...]
#
# Examples:
#   ./scripts/run-isolated-instance.sh agent-a gateway
#   ./scripts/run-isolated-instance.sh agent-b launcher -- -console -no-browser
#   PICOCLAW_GATEWAY_PORT=18010 ./scripts/run-isolated-instance.sh agent-a gateway
#   PICOCLAW_LAUNCHER_PORT=18810 ./scripts/run-isolated-instance.sh agent-a launcher
#
# Optional environment:
#   PICOCLAW_INSTANCE_ROOT       Base directory for instances.
#                                Default: $HOME/.picoclaw-instances
#   PICOCLAW_INSTANCE_HOME       Exact home directory for this instance.
#                                Default: $PICOCLAW_INSTANCE_ROOT/<name>
#   PICOCLAW_BINARY              picoclaw CLI path. Default: picoclaw
#   PICOCLAW_WEB_BINARY          launcher/web binary path. Default: picoclaw-web
#   PICOCLAW_GATEWAY_PORT        Gateway port override.
#   PICOCLAW_LAUNCHER_PORT       Launcher port override.
#   PICOCLAW_GATEWAY_HOST        Gateway bind host. Default: localhost
#   PICOCLAW_LAUNCHER_HOST       Launcher bind host, if needed.

set -eu

usage() {
    sed -n '2,35p' "$0" >&2
    exit 2
}

if [ "$#" -lt 1 ]; then
    usage
fi

INSTANCE_NAME=$1
shift

case "$INSTANCE_NAME" in
    ""|.*|*/*|*' '*|*':'*)
        echo "ERROR: instance name must be a simple path-safe name, for example agent-a" >&2
        exit 2
        ;;
esac

MODE=${1:-gateway}
if [ "$#" -gt 0 ]; then
    shift
fi

if [ "${1:-}" = "--" ]; then
    shift
fi

case "$MODE" in
    gateway|launcher) ;;
    *)
        echo "ERROR: mode must be gateway or launcher" >&2
        usage
        ;;
esac

if [ -z "${HOME:-}" ]; then
    echo "ERROR: HOME is not set" >&2
    exit 1
fi

INSTANCE_ROOT=${PICOCLAW_INSTANCE_ROOT:-"$HOME/.picoclaw-instances"}
INSTANCE_HOME=${PICOCLAW_INSTANCE_HOME:-"$INSTANCE_ROOT/$INSTANCE_NAME"}

# Stable per-name port offset. Collisions are possible but uncommon; override
# PICOCLAW_GATEWAY_PORT / PICOCLAW_LAUNCHER_PORT when you need exact ports.
PORT_OFFSET=$(printf '%s' "$INSTANCE_NAME" | cksum | awk '{ print ($1 % 1000) + 1 }')
DEFAULT_GATEWAY_PORT=$((18000 + PORT_OFFSET))
DEFAULT_LAUNCHER_PORT=$((18800 + PORT_OFFSET))

GATEWAY_PORT=${PICOCLAW_GATEWAY_PORT:-$DEFAULT_GATEWAY_PORT}
LAUNCHER_PORT=${PICOCLAW_LAUNCHER_PORT:-$DEFAULT_LAUNCHER_PORT}
GATEWAY_HOST=${PICOCLAW_GATEWAY_HOST:-localhost}

PICOC_LI=${PICOCLAW_BINARY:-picoclaw}
PICOC_WEB=${PICOCLAW_WEB_BINARY:-picoclaw-web}

mkdir -p \
    "$INSTANCE_HOME" \
    "$INSTANCE_HOME/workspace" \
    "$INSTANCE_HOME/logs" \
    "$INSTANCE_HOME/skills" \
    "$INSTANCE_HOME/cache" \
    "$INSTANCE_HOME/state"

CONFIG_PATH=$INSTANCE_HOME/config.json

export PICOCLAW_HOME=$INSTANCE_HOME
export PICOCLAW_CONFIG=$CONFIG_PATH
export PICOCLAW_GATEWAY_HOST=$GATEWAY_HOST
export PICOCLAW_GATEWAY_PORT=$GATEWAY_PORT
export PICOCLAW_AGENTS_DEFAULTS_WORKSPACE=$INSTANCE_HOME/workspace

echo "PicoClaw isolated instance"
echo "  name:          $INSTANCE_NAME"
echo "  home:          $PICOCLAW_HOME"
echo "  config:        $PICOCLAW_CONFIG"
echo "  workspace:     $PICOCLAW_AGENTS_DEFAULTS_WORKSPACE"
echo "  gateway:       $PICOCLAW_GATEWAY_HOST:$PICOCLAW_GATEWAY_PORT"
if [ "$MODE" = "launcher" ]; then
    echo "  launcher port: $LAUNCHER_PORT"
fi
echo ""

case "$MODE" in
    gateway)
        exec "$PICOC_LI" gateway "$@"
        ;;
    launcher)
        exec "$PICOC_WEB" -port "$LAUNCHER_PORT" "$@" "$CONFIG_PATH"
        ;;
esac
