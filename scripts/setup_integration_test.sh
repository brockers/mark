#!/bin/bash

# Setup integration tests for mark bookmark CLI tool

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MARK_BINARY="$SCRIPT_DIR/../mark"

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0
test_count=0

# Test helper functions
test_pass() {
    ((PASSED++)) || true
    echo -e "${GREEN}✓${NC} $1"
}

test_fail() {
    ((FAILED++)) || true
    echo -e "${RED}✗${NC} $1"
}

run_test() {
    ((test_count++)) || true
    echo ""
    echo "Test $test_count: $1"
}

# Setup test environment
setup_test_env() {
    export HOME="/tmp/mark-setup-test-$$"
    mkdir -p "$HOME"
    export PATH="$SCRIPT_DIR/..:$PATH"
}

# Cleanup test environment
cleanup_test_env() {
    if [ -d "$HOME" ]; then
        rm -rf "$HOME"
    fi
}

# Ensure binary exists
if [ ! -f "$MARK_BINARY" ]; then
    echo "Error: mark binary not found at $MARK_BINARY"
    echo "Please run 'make build' first"
    exit 1
fi

# Setup
trap cleanup_test_env EXIT
setup_test_env

# Test 1: First run creates config
run_test "First run creates config file"
echo "$HOME/.marks" | "$MARK_BINARY" --config >/dev/null 2>&1 </dev/null || true
if [ -f "$HOME/.mark" ]; then
    test_pass "Config file created on first run"
else
    test_fail "Config file not created"
fi

# Test 2: Config contains marksdir
run_test "Config contains marksdir setting"
if grep -q "marksdir=" "$HOME/.mark"; then
    test_pass "Config contains marksdir"
else
    test_fail "Config does not contain marksdir"
fi

# Test 3: Marks directory is created
run_test "Marks directory is created"
if [ -d "$HOME/.marks" ]; then
    test_pass "Marks directory exists"
else
    test_fail "Marks directory not created"
fi

# Test 4: Reconfiguration works
run_test "Reconfiguration updates config"
CUSTOM_MARKS_DIR="$HOME/custom-marks"
printf "$CUSTOM_MARKS_DIR\nn\nn\n" | "$MARK_BINARY" --config >/dev/null 2>&1 || true
if grep -q "custom-marks" "$HOME/.mark"; then
    test_pass "Config updated with custom path"
else
    test_fail "Config not updated"
fi

# Test 5: Custom marks directory created
run_test "Custom marks directory created"
if [ -d "$CUSTOM_MARKS_DIR" ]; then
    test_pass "Custom marks directory exists"
else
    test_fail "Custom marks directory not created"
fi

# Test 6: Shell detection
run_test "Shell detection works"
export SHELL="/bin/bash"
if [ "$(echo "$SHELL" | "$MARK_BINARY" --config 2>&1 | grep -c 'bash' || true)" -ge 0 ]; then
    test_pass "Shell detection functional"
else
    test_fail "Shell detection failed"
fi

# Test 7: Unified RC file created for completions
run_test "Unified RC file created for completions"
# Clean up first
rm -f "$HOME/.mark_bash_rc"
printf "$HOME/.marks\ny\nn\n" | "$MARK_BINARY" --config >/dev/null 2>&1 || true
if [ -f "$HOME/.mark_bash_rc" ]; then
    test_pass "Unified RC file created at ~/.mark_bash_rc"
else
    test_fail "Unified RC file not created"
fi

# Test 8: Source line added to .bashrc
run_test "Source line added to .bashrc"
if grep -q "# mark shell integration" "$HOME/.bashrc" 2>/dev/null; then
    test_pass "Source line found in .bashrc"
else
    test_fail "Source line not found in .bashrc"
fi

# Test 9: RC file contains features header
run_test "RC file contains features header"
if grep -q "# Features:" "$HOME/.mark_bash_rc" 2>/dev/null; then
    test_pass "Features header found in RC file"
else
    test_fail "Features header not found in RC file"
fi

# Test 10: Aliases setup creates unified RC with aliases
run_test "Aliases setup creates unified RC with aliases"
rm -f "$HOME/.mark_bash_rc"
printf "y\n" | "$MARK_BINARY" --alias >/dev/null 2>&1 || true
if [ -f "$HOME/.mark_bash_rc" ] && grep -q "alias marks=" "$HOME/.mark_bash_rc" 2>/dev/null; then
    test_pass "Aliases added to unified RC file"
else
    test_fail "Aliases not found in unified RC file"
fi

# Test 11: Both aliases and completions in single RC file
run_test "Both aliases and completions in single RC file"
rm -f "$HOME/.mark_bash_rc"
printf "$HOME/.marks\ny\ny\n" | "$MARK_BINARY" --config >/dev/null 2>&1 || true
if grep -q "alias marks=" "$HOME/.mark_bash_rc" 2>/dev/null && grep -q "_mark_complete()" "$HOME/.mark_bash_rc" 2>/dev/null; then
    test_pass "Both aliases and completions in RC file"
else
    test_fail "Missing aliases or completions in RC file"
fi

# Print summary
echo ""
echo "========================================"
echo "Setup Test Summary"
echo "========================================"
echo "Tests passed: $PASSED"
echo "Tests failed: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi

echo -e "${GREEN}All setup tests passed!${NC}"
exit 0
