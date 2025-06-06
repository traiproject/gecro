# Path to the commit message file
COMMIT_MSG_FILE="$1"
COMMIT_MSG=$(<"$COMMIT_MSG_FILE")

# Conventional Commits regex:
#   <type>[optional scope]: <description>
#   Supported types: build, ci, chore, docs, feat, fix, perf, refactor, revert, style, test
CONV_REGEX="^(build|ci|chore|docs|feat|fix|perf|refactor|revert|style|test)(\([[:alnum:]\-]+\))?:\s.+$"

# 1. Validate header format
if ! echo "$COMMIT_MSG" | grep -qE "$CONV_REGEX"; then
  echo "Commit message does not follow Conventional Commits format."
  echo "Format: <type>(<scope>): <description>"
  echo "Example: feat(parser): add support for nested arrays"
  exit 1
fi

# 2. Enforce blank line before body (if a body exists)
#    If the commit message has more than 2 lines, then there must be
#    exactly one blank line (i.e. an empty second line) separating header and body.
TOTAL_LINES=$(wc -l < "$COMMIT_MSG_FILE")
if [ "$TOTAL_LINES" -gt 2 ]; then
  # Read the second line
  SECOND_LINE=$(sed -n '2p' "$COMMIT_MSG_FILE")
  if [ -n "$SECOND_LINE" ]; then
    echo "Commit message body must be separated from header by a blank line."
    echo "Ensure you have an empty line (no spaces) after the header before starting the body."
    exit 1
  fi
fi

# 3. If everything’s OK, exit 0
exit 0
