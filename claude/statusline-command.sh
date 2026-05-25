#!/usr/bin/env bash

input=$(cat)

model=$(echo "$input" | jq -r '.model.display_name // "unknown"')
used=$(echo "$input" | jq -r '.context_window.used_percentage // empty')
used_tokens=$(echo "$input" | jq -r '.context_window.total_input_tokens // empty')
max_tokens=$(echo "$input" | jq -r '.context_window.context_window_size // empty')
effort=$(echo "$input" | jq -r '.effort.level // empty')

# git branch + short cwd
cwd=$(echo "$input" | jq -r '.cwd // empty')
branch=""
short_cwd=""
if [ -n "$cwd" ]; then
  branch=$(git -C "$cwd" --no-optional-locks rev-parse --abbrev-ref HEAD 2>/dev/null)
  tilde_cwd=$(echo "$cwd" | sed "s|^$HOME|~|")
  depth=$(echo "$tilde_cwd" | tr -cd '/' | wc -c)
  if [ "$depth" -gt 1 ]; then
    short_cwd="~/…/$(basename "$cwd")"
  else
    short_cwd="$tilde_cwd"
  fi
fi

# build parts
parts=""

# model
parts="$model"

# cwd
if [ -n "$short_cwd" ]; then
  parts="$parts  |  $short_cwd"
fi

# context usage
if [ -n "$used" ]; then
  printf -v used_fmt "%.0f" "$used"
  ctx_str="uso: ${used_fmt}%"
  if [ -n "$used_tokens" ] && [ -n "$max_tokens" ]; then
    used_k=$(awk "BEGIN {printf \"%.0f\", $used_tokens/1000}")
    max_k=$(awk "BEGIN {printf \"%.0f\", $max_tokens/1000}")
    ctx_str="${ctx_str} · ${used_k}k/${max_k}k"
  fi
  parts="$parts  |  $ctx_str"
fi

# git branch
if [ -n "$branch" ] && [ "$branch" != "HEAD" ]; then
  parts="$parts  |  $branch"
fi

# effort
if [ -n "$effort" ]; then
  parts="$parts  |  effort: $effort"
fi

printf "%s" "$parts"
