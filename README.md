# prompto

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/krostar/prompto)
[![Licence](https://img.shields.io/github/license/krostar/prompto.svg?style=for-the-badge)](https://tldrlegal.com/license/gnu-lesser-general-public-license-v3-(lgpl-3))
![Latest version](https://img.shields.io/github/tag/krostar/prompto.svg?style=for-the-badge)

A fast, modulable and extremely configurable prompt for any shells.

## Motivation

// TODO

## Installation

// TODO

### Bash

```bash
_prompt_command_timer_start() {
    { [ ! -t 1 ] || [ -p /dev/stdout ] || [[ ! -t 1 && ! -p /dev/stdout ]] || [ -n "$COMP_LINE" ] || [[ "$BASH_COMMAND" == "$PROMPT_COMMAND" ]]; } && return
    prompt_command_started_at=$(date +'%s')
}

_prompt() {
    local -r last_cmd_status="$?"
    local last_cmd_duration=0
    
    if [ -n "$prompt_command_started_at" ]; then
        last_cmd_duration="$((($(date +'%s') - prompt_command_started_at) * 1000000000))"
        unset prompt_command_started_at
    fi

    local -r LPROMPT="$(prompto --left -s $last_cmd_status -d $last_cmd_duration)"
    local -r RPROMPT="$(prompto --right)"
    local -r RPROMPT_NOCOLOR="$(echo "$RPROMPT" | sed $'s,\x1b\\[[0-9;]*[a-zA-Z],,g')"
    local -r RPROMPT_LEN="$((${#RPROMPT_NOCOLOR}-1))"

    PS1="\[\e[s\e[${COLUMNS:-$(tput cols)}C\e[${RPROMPT_LEN}D${RPROMPT}\e[u\]${LPROMPT}"
}

PROMPT_COMMAND=_prompt
trap _prompt_command_timer_start DEBUG
```

### ZSH

```zsh
zmodload zsh/datetime
zmodload zsh/mathfunc

function preexec() {
    prompt_command_started_at="$EPOCHREALTIME"
}

function _prompt() {
    local last_cmd_status=$?
    local last_cmd_duration=0

    if [ -n "$prompt_command_started_at" ]; then
        last_cmd_duration="$((int(rint((EPOCHREALTIME - prompt_command_started_at) * 1000000000))))"
        unset prompt_command_started_at
    fi

    PROMPT="$(prompto/build/bin/prompto --left --shell zsh -s $last_cmd_status -d $last_cmd_duration)"
    RPROMPT="$(prompto/build/bin/prompto --right --shell zsh)"
    ZLE_RPROMPT_INDENT=0

    export PROMPT RPROMPT ZLE_RPROMPT_INDENT
}

function _init_prompt() {
    for s in "${precmd_functions[@]}"; do
        if [ "$s" = "_prompt" ]; then
            return
        fi
    done
    precmd_functions+=(_prompt)
}

if [ "$TERM" != "linux" ]; then
    _init_prompt
fi
```

### Fish

```fish
function fish_prompt
    set -l exitcode $status
    set -l duration (math "$CMD_DURATION * 1000000")
    prompto --left -d $duration -s $exitcode
end

function fish_right_prompt
    prompto --right
end
```

## Configuration

// TODO

## Limitation

// TODO

## Comparaison

// TODO
