# prompto

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/krostar/prompto)
[![Licence](https://img.shields.io/github/license/krostar/prompto.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/prompto.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/prompto/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/prompto)
[![Code quality](https://img.shields.io/codacy/grade/xxxxxx/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/prompto/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/xxxxxx.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/prompto/dashboard)

A fast, modulable and extremely configurable prompt for any shells.

## Installation

**prompto** has been tested on linux and macOS, and is available using:

-   the [releases assets of github](https://github.com/krostar/prompto/releases)
-   package manager: 
    -   ~~apt _(available soon)_~~
    -   ~~brew _(available soon)_~~
    -   ~~aur _(available soon)_~~

## Configuration

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

    PROMPT="$(/Users/alexis.destrez/Work/Perso/go/prompto/build/bin/prompto --left --shell zsh -s $last_cmd_status -d $last_cmd_duration)"
    RPROMPT="$(/Users/alexis.destrez/Work/Perso/go/prompto/build/bin/prompto --right --shell zsh)"
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
    set duration (math "$CMD_DURATION * 1000000")
    $HOME/Work/Perso/go/prompto/build/bin/prompto --left -d $duration -s $exitcode
end

function fish_right_prompt
    $HOME/Work/Perso/go/prompto/build/bin/prompto --right
end
```
