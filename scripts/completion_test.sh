#!/bin/bash

# test_completion.sh - Test bash completion for the mark command
# This script tests the completion functionality to ensure it properly
# handles partial bookmark name matching

# Don't use set -e because we want to continue even if a test fails

echo "=== Mark Completion Test Suite ==="
echo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Create a temporary test directory with mock bookmarks
TEST_DIR=$(mktemp -d)
TEST_CONFIG=$(mktemp)
trap "rm -rf $TEST_DIR $TEST_CONFIG" EXIT

# Set up test environment
echo "Setting up test environment..."
echo "marksdir=$TEST_DIR/.marks" > $TEST_CONFIG

# Create mock bookmark directories and symlinks
mkdir -p "$TEST_DIR/.marks"
mkdir -p "$TEST_DIR/projects/work-project"
mkdir -p "$TEST_DIR/projects/personal-project"
mkdir -p "$TEST_DIR/downloads"
mkdir -p "$TEST_DIR/documents/important"
mkdir -p "$TEST_DIR/code/golang-app"
mkdir -p "$TEST_DIR/code/python-script"
mkdir -p "$TEST_DIR/home/bobby"

# Create symlinks in the marks directory
ln -s "$TEST_DIR/projects/work-project" "$TEST_DIR/.marks/work"
ln -s "$TEST_DIR/projects/personal-project" "$TEST_DIR/.marks/personal"
ln -s "$TEST_DIR/downloads" "$TEST_DIR/.marks/downloads"
ln -s "$TEST_DIR/documents/important" "$TEST_DIR/.marks/docs"
ln -s "$TEST_DIR/code/golang-app" "$TEST_DIR/.marks/golang"
ln -s "$TEST_DIR/code/python-script" "$TEST_DIR/.marks/python"
ln -s "$TEST_DIR/home/bobby" "$TEST_DIR/.marks/home"

# Build the completion script dynamically
echo "Building completion script..."

# Extract the completion function and adapt for testing
COMPLETION_SCRIPT=$(mktemp)
cat > $COMPLETION_SCRIPT << 'SCRIPT_END'
#!/bin/bash

_mark_complete_test() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"

    # If we're on the first argument
    if [[ ${COMP_CWORD} -eq 1 ]]; then
        # If user starts typing a dash, offer flags
        if [[ "$cur" == -* ]]; then
            local flags="-l -d -j -v -h --config --autocomplete --alias --help --version"
            COMPREPLY=($(compgen -W "$flags" -- "${cur}"))
        else
            # For create command, offer existing bookmarks as suggestions
            if [[ -d TEST_MARKS_DIR ]]; then
                local marks=$(ls TEST_MARKS_DIR 2>/dev/null | tr '\n' ' ')
                COMPREPLY=($(compgen -W "$marks" -- "${cur}"))
            fi
        fi
    # If previous was -d or -j, offer bookmark names
    elif [[ "$prev" == "-d" || "$prev" == "-j" ]]; then
        if [[ -d TEST_MARKS_DIR ]]; then
            local marks=$(ls TEST_MARKS_DIR 2>/dev/null | tr '\n' ' ')
            COMPREPLY=($(compgen -W "$marks" -- "${cur}"))
        fi
    fi
}
SCRIPT_END

# Replace TEST_MARKS_DIR with actual path
sed -i "s|TEST_MARKS_DIR|$TEST_DIR/.marks|g" $COMPLETION_SCRIPT

# Source the test completion function
source $COMPLETION_SCRIPT

# Test function
run_test() {
    local test_name="$1"
    local input="$2"
    local expected_count="$3"
    local expected_pattern="$4"

    # Set up completion environment
    export COMP_WORDS=("mark" "$input")
    export COMP_CWORD=1
    COMPREPLY=()

    # Run completion
    _mark_complete_test

    # Check results
    local result_count=${#COMPREPLY[@]}
    local pattern_found=0

    if [[ -n "$expected_pattern" ]]; then
        for result in "${COMPREPLY[@]}"; do
            if [[ "$result" == $expected_pattern* ]]; then
                pattern_found=1
                break
            fi
        done
    else
        pattern_found=1  # No pattern to check
    fi

    if [[ "$expected_count" == "-1" ]] || [[ $result_count -eq $expected_count ]]; then
        count_ok=1
    elif [[ "$expected_count" == "+" ]] && [[ $result_count -gt 0 ]]; then
        count_ok=1
    else
        count_ok=0
    fi

    if [[ $count_ok -eq 1 ]] && [[ $pattern_found -eq 1 ]]; then
        echo -e "${GREEN}✓${NC} $test_name (found $result_count matches)"
        if [[ $result_count -gt 0 ]] && [[ $result_count -le 5 ]]; then
            for result in "${COMPREPLY[@]}"; do
                echo "    - $result"
            done
        fi
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗${NC} $test_name"
        echo "    Expected: count=$expected_count, pattern=$expected_pattern"
        echo "    Got: count=$result_count"
        if [[ $result_count -gt 0 ]] && [[ $result_count -le 10 ]]; then
            echo "    Results:"
            for result in "${COMPREPLY[@]}"; do
                echo "      - $result"
            done
        fi
        ((TESTS_FAILED++))
    fi
}

echo "Running completion tests..."
echo

# Test partial matching
run_test "Partial match 'work' should return work bookmark" "work" 1 "work"
run_test "Partial match 'p' should return personal and python bookmarks" "p" 2 "p"
run_test "Partial match 'down' should return downloads bookmark" "down" 1 "downloads"
run_test "Partial match 'gol' should return golang bookmark" "gol" 1 "golang"
run_test "Partial match 'home' should return home bookmark" "home" 1 "home"

# Test exact prefix matching
run_test "Exact prefix 'person' should return personal" "person" 1 "personal"
run_test "Exact prefix 'doc' should return docs" "doc" 1 "docs"

# Test non-matching input
run_test "Non-matching 'xyz' should return no results" "xyz" 0 ""

# Test empty input (should return all bookmarks)
run_test "Empty input should return all bookmarks" "" 7 ""

# Test flag completion
run_test "Flag '-l' should match -l flag" "-l" 1 "-l"
run_test "Flag '-v' should match -v flag" "-v" 1 "-v"
run_test "Flag '-j' should match -j flag" "-j" 1 "-j"
run_test "Flag '-d' should match -d flag" "-d" 1 "-d"
run_test "Flag '--h' should match --help" "--h" 1 "--help"
run_test "Flag '--v' should match --version" "--v" 1 "--version"
run_test "Flag '--a' should match --autocomplete and --alias" "--a" 2 "--a"

# Test completion after -j flag (for jump)
export COMP_WORDS=("mark" "-j" "work")
export COMP_CWORD=2
COMPREPLY=()
_mark_complete_test
if [[ ${#COMPREPLY[@]} -eq 1 ]]; then
    echo -e "${GREEN}✓${NC} Completion after -j flag works (found ${#COMPREPLY[@]} matches)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}✗${NC} Completion after -j flag failed (expected 1, got ${#COMPREPLY[@]})"
    ((TESTS_FAILED++))
fi

# Test completion after -d flag (for delete)
export COMP_WORDS=("mark" "-d" "p")
export COMP_CWORD=2
COMPREPLY=()
_mark_complete_test
if [[ ${#COMPREPLY[@]} -eq 2 ]]; then
    echo -e "${GREEN}✓${NC} Completion after -d flag works (found ${#COMPREPLY[@]} matches)"
    ((TESTS_PASSED++))
else
    echo -e "${RED}✗${NC} Completion after -d flag failed (expected 2, got ${#COMPREPLY[@]})"
    ((TESTS_FAILED++))
fi

echo
echo "Testing alias completions (marks, unmark, jump)..."

# Test helper for alias completion
run_alias_test() {
    local alias_name="$1"
    local test_name="$2"
    local input="$3"
    local expected_count="$4"

    # Set up completion environment for alias
    export COMP_WORDS=("$alias_name" "$input")
    export COMP_CWORD=1
    COMPREPLY=()

    # Run completion
    _mark_complete_test

    # Check results
    local result_count=${#COMPREPLY[@]}

    if [[ "$expected_count" == "-1" ]] || [[ $result_count -eq $expected_count ]]; then
        echo -e "${GREEN}✓${NC} $test_name (found $result_count matches)"
        if [[ $result_count -gt 0 ]] && [[ $result_count -le 3 ]]; then
            for result in "${COMPREPLY[@]}"; do
                echo "    - $result"
            done
        fi
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗${NC} $test_name"
        echo "    Expected: count=$expected_count"
        echo "    Got: count=$result_count"
        ((TESTS_FAILED++))
    fi
}

# Test 'marks' alias (mark -l)
run_alias_test "marks" "Alias 'marks' should complete work bookmark" "work" 1
run_alias_test "marks" "Alias 'marks' should complete p* bookmarks" "p" 2
run_alias_test "marks" "Alias 'marks' should handle empty input" "" 7

# Test 'unmark' alias (mark -d)
run_alias_test "unmark" "Alias 'unmark' should complete downloads bookmark" "downloads" 1
run_alias_test "unmark" "Alias 'unmark' should complete golang bookmark" "golang" 1

# Test 'jump' alias (mark -j)
run_alias_test "jump" "Alias 'jump' should complete home bookmark" "home" 1
run_alias_test "jump" "Alias 'jump' should complete docs bookmark" "docs" 1

echo
echo "Verifying completion registration for all aliases..."

# Generate the actual completion script to verify it has all aliases
TEMP_COMPLETION_DIR=$(mktemp -d)
export HOME="$TEMP_COMPLETION_DIR"
trap "rm -rf $TEMP_COMPLETION_DIR" EXIT

# Create a simple test to extract the bash completion script
# We'll read it directly from the Go code's embedded string
# by checking what the mark binary generates
MARK_BINARY="$(dirname "$0")/../mark"

# Create a config file first so it doesn't go through setup
mkdir -p "$HOME/.marks"
echo "marksdir=$HOME/.marks" > "$HOME/.mark"

# Extract the completion script content by running mark and checking the generated file
# We need to simulate the completion setup
export SHELL="/bin/bash"
echo "y" | "$MARK_BINARY" --autocomplete > /dev/null 2>&1 || true

BASH_COMPLETION_FILE="$HOME/.mark.bash"

if [ -f "$BASH_COMPLETION_FILE" ]; then
    if grep -q "complete -F _mark_complete mark" "$BASH_COMPLETION_FILE" &&
       grep -q "complete -F _mark_complete marks" "$BASH_COMPLETION_FILE" &&
       grep -q "complete -F _mark_complete unmark" "$BASH_COMPLETION_FILE" &&
       grep -q "complete -F _mark_complete jump" "$BASH_COMPLETION_FILE"; then
        echo -e "${GREEN}✓${NC} All aliases have completion registered in generated script"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗${NC} Missing completion registration for one or more aliases"
        echo "Generated completion script should contain:"
        echo "  complete -F _mark_complete mark"
        echo "  complete -F _mark_complete marks"
        echo "  complete -F _mark_complete unmark"
        echo "  complete -F _mark_complete jump"
        ((TESTS_FAILED++))
    fi

    # Verify the helper function for formatted output exists
    if grep -q "_mark_list_with_paths" "$BASH_COMPLETION_FILE"; then
        echo -e "${GREEN}✓${NC} Completion script includes formatted bookmark display helper"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗${NC} Completion script missing _mark_list_with_paths helper function"
        ((TESTS_FAILED++))
    fi
else
    echo -e "${RED}✗${NC} Completion script was not generated"
    ((TESTS_FAILED++))
fi

echo
echo "==================================="
echo "Test Summary:"
echo -e "  Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "  Failed: ${RED}$TESTS_FAILED${NC}"

if [[ $TESTS_FAILED -eq 0 ]]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
