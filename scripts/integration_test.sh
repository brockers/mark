#!/bin/bash

# Integration tests for mark bookmark CLI tool

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MARK_BINARY="$SCRIPT_DIR/../mark"

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0

# Test counter
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
    export HOME="/tmp/mark-test-$$"
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

# Test 1: Version flag
run_test "Version flag works"
if "$MARK_BINARY" --version | grep -q "dev"; then
    test_pass "Version flag returns version"
else
    test_fail "Version flag did not return version"
fi

# Test 2: Help flag
run_test "Help flag works"
if "$MARK_BINARY" --help | grep -q "mark - A minimalist CLI bookmark tool"; then
    test_pass "Help flag displays help text"
else
    test_fail "Help flag did not display help text"
fi

# Test 3: Create default bookmark
run_test "Create default bookmark"
TEST_DIR="$HOME/test-project"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Non-interactive config creation
printf "$HOME/.marks\n\n\n" | "$MARK_BINARY" --config >/dev/null 2>&1 || true

# Create bookmark
if "$MARK_BINARY" 2>/dev/null | grep -q "Created bookmark"; then
    test_pass "Created default bookmark"
else
    test_fail "Failed to create default bookmark"
fi

# Test 4: List bookmarks
run_test "List bookmarks"
if "$MARK_BINARY" -l 2>/dev/null | grep -q "test-project"; then
    test_pass "List shows created bookmark"
else
    test_fail "List did not show bookmark"
fi

# Test 5: Create named bookmark
run_test "Create named bookmark"
TEST_DIR2="$HOME/another-project"
mkdir -p "$TEST_DIR2"
cd "$TEST_DIR2"

if "$MARK_BINARY" mymark 2>/dev/null | grep -q "Created bookmark 'mymark'"; then
    test_pass "Created named bookmark"
else
    test_fail "Failed to create named bookmark"
fi

# Test 6: Jump to bookmark (path output)
run_test "Jump to bookmark"
JUMP_OUTPUT=$("$MARK_BINARY" -j mymark 2>/dev/null)
if [ "$JUMP_OUTPUT" = "$TEST_DIR2" ]; then
    test_pass "Jump returned correct path"
else
    test_fail "Jump did not return correct path (got: $JUMP_OUTPUT)"
fi

# Test 7: Delete bookmark
run_test "Delete bookmark"
if "$MARK_BINARY" -d mymark 2>/dev/null | grep -q "Removed bookmark 'mymark'"; then
    test_pass "Deleted bookmark"
else
    test_fail "Failed to delete bookmark"
fi

# Test 8: List after delete
run_test "List after delete shows bookmark is gone"
if ! "$MARK_BINARY" -l 2>/dev/null | grep -q "mymark"; then
    test_pass "Deleted bookmark no longer appears in list"
else
    test_fail "Deleted bookmark still appears in list"
fi

# Test 9: Broken symlink detection
run_test "Broken symlink detection"
# Create a bookmark, then delete the target directory
BROKEN_DIR="$HOME/will-be-deleted"
mkdir -p "$BROKEN_DIR"
cd "$BROKEN_DIR"
"$MARK_BINARY" brokenmark >/dev/null 2>&1
rm -rf "$BROKEN_DIR"

if "$MARK_BINARY" -l 2>/dev/null | grep "brokenmark" | grep -q "\[broken\]"; then
    test_pass "Broken symlink detected and marked"
else
    test_fail "Broken symlink not properly detected"
fi

# Print summary
echo ""
echo "========================================"
echo "Integration Test Summary"
echo "========================================"
echo "Tests passed: $PASSED"
echo "Tests failed: $FAILED"
echo "========================================"

if [ $FAILED -gt 0 ]; then
    exit 1
fi

echo -e "${GREEN}All integration tests passed!${NC}"
exit 0
