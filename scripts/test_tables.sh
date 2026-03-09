#!/usr/bin/env bash
# Run a smoke-test query against every torii table and report results.
# Exit code is the number of failed tables (0 = all passed).

set -euo pipefail

# Colour codes (disabled if not a terminal)
if [ -t 1 ]; then
  GREEN="\033[0;32m"
  RED="\033[0;31m"
  YELLOW="\033[0;33m"
  RESET="\033[0m"
else
  GREEN="" RED="" YELLOW="" RESET=""
fi

PASS=0
FAIL=0
SKIP=0

run_test() {
  local table="$1"
  local query="$2"
  printf "  %-40s" "$table"
  local output
  if output=$(steampipe query "$query" 2>&1); then
    printf "${GREEN}PASS${RESET}\n"
    ((PASS++)) || true
  else
    # A 403 / missing scope is treated as a skip, not a hard failure
    if echo "$output" | grep -q "missing_required_scope\|403\|Forbidden"; then
      printf "${YELLOW}SKIP${RESET} (insufficient permissions)\n"
      ((SKIP++)) || true
    else
      printf "${RED}FAIL${RESET}\n"
      echo "$output" | sed 's/^/    /'
      ((FAIL++)) || true
    fi
  fi
}

echo ""
echo "Torii Steampipe plugin — table smoke tests"
echo "==========================================="
echo ""

run_test "torii_user"       "select id, email, lifecycle_status from torii_user limit 1"
run_test "torii_app"        "select id, name, state from torii_app limit 1"
run_test "torii_app_user"   "select app_id, id_user, email, status from torii_app_user where app_id = 1 limit 1"
run_test "torii_contract"   "select id, name, status from torii_contract limit 1"
run_test "torii_role"       "select id, name, is_admin from torii_role limit 1"
run_test "torii_audit_log"  "select performed_by_email, type, creation_time from torii_audit_log limit 1"

echo ""
printf "Results: ${GREEN}%d passed${RESET}, ${RED}%d failed${RESET}, ${YELLOW}%d skipped${RESET}\n" "$PASS" "$FAIL" "$SKIP"
echo ""

exit "$FAIL"
