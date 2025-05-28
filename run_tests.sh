#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print error and exit
error() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

# Function to print success
success() {
    echo -e "${GREEN}$1${NC}"
}

# Function to print warning
warning() {
    echo -e "${YELLOW}$1${NC}"
}

# Help message
show_help() {
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -i, --integration-only  Run only integration tests"
    echo "  -k, --api-key KEY       Set the Vartiq API key"
    echo "  -u, --api-url URL       Set the Vartiq API URL (optional)"
    echo "  -v, --verbose           Enable verbose output"
    echo
    echo "Environment variables:"
    echo "  VARTIQ_API_KEY          Vartiq API key (can be used instead of -k)"
    echo "  VARTIQ_API_URL          Vartiq API URL (can be used instead of -u)"
}

# Default values
INTEGRATION_ONLY=false
VERBOSE=false
API_KEY=${VARTIQ_API_KEY:-""}
API_URL=${VARTIQ_API_URL:-""}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -i|--integration-only)
            INTEGRATION_ONLY=true
            shift
            ;;
        -k|--api-key)
            API_KEY="$2"
            shift 2
            ;;
        -u|--api-url)
            API_URL="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Check if API key is provided
if [ -z "$API_KEY" ]; then
    error "API key is required. Set it using -k/--api-key option or VARTIQ_API_KEY environment variable"
fi

# Export environment variables
export VARTIQ_API_KEY="$API_KEY"
if [ -n "$API_URL" ]; then
    export VARTIQ_API_URL="$API_URL"
fi

# Build test command
TEST_CMD="go test ./vartiq/..."
if [ "$INTEGRATION_ONLY" = true ]; then
    TEST_CMD="go test ./vartiq/... -run Integration"
fi
if [ "$VERBOSE" = true ]; then
    TEST_CMD="$TEST_CMD -v"
fi

# Run tests
echo "Running tests..."
if $TEST_CMD; then
    success "\nAll tests passed successfully! 🎉"
else
    error "Some tests failed"
fi 

# Usage
# ./run_tests.sh -k your_api_key -u https://url.com