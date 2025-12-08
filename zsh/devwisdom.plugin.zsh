#!/usr/bin/env zsh
# devwisdom zsh plugin
# Provides convenient functions for accessing wisdom quotes and advisor consultations

# Check if devwisdom is available
if ! command -v devwisdom &> /dev/null; then
    echo "devwisdom: command not found. Please install devwisdom first."
    return 1
fi

# Get the path to devwisdom binary
DEVWISDOM_CMD="${DEVWISDOM_CMD:-devwisdom}"

# Show daily wisdom (can be called on shell startup)
devwisdom-daily() {
    local source="${1:-}"
    if [[ -n "$source" ]]; then
        $DEVWISDOM_CMD quote --source "$source"
    else
        $DEVWISDOM_CMD quote
    fi
}

# Get a quick quote
devwisdom-quote() {
    local source="${1:-}"
    local score="${2:-50}"
    
    if [[ -n "$source" ]]; then
        $DEVWISDOM_CMD quote --source "$source" --score "$score"
    else
        $DEVWISDOM_CMD quote --score "$score"
    fi
}

# Consult an advisor
devwisdom-consult() {
    local metric="${1:-}"
    local tool="${2:-}"
    local stage="${3:-}"
    local score="${4:-50}"
    
    local args=()
    [[ -n "$metric" ]] && args+=("--metric" "$metric")
    [[ -n "$tool" ]] && args+=("--tool" "$tool")
    [[ -n "$stage" ]] && args+=("--stage" "$stage")
    args+=("--score" "$score")
    
    $DEVWISDOM_CMD consult "${args[@]}"
}

# List available sources
devwisdom-sources() {
    $DEVWISDOM_CMD sources "$@"
}

# List available advisors
devwisdom-advisors() {
    $DEVWISDOM_CMD advisors "$@"
}

# Get daily briefing
devwisdom-briefing() {
    local days="${1:-1}"
    $DEVWISDOM_CMD briefing --days "$days"
}

# Auto-show daily wisdom on shell startup (if enabled)
if [[ "${DEVWISDOM_AUTO_DAILY:-false}" == "true" ]]; then
    devwisdom-daily
fi
