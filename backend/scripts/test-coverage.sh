#!/bin/bash
# test-coverage.sh - Run comprehensive test coverage analysis for the importer package

set -e

echo "==================================="
echo "BaMoRT Import Package Test Coverage"
echo "==================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Change to backend directory
cd "$(dirname "$0")"

# Create coverage directory
mkdir -p coverage

echo "Running unit tests with coverage..."
go test -v -coverprofile=coverage/unit.out ./importer/

echo ""
echo "Running integration tests..."
go test -v -tags=integration ./importer/ -run "Test.*_test"

echo ""
echo "Running E2E tests..."
go test -v -tags=e2e ./importer/e2e_test.go -timeout 30s

echo ""
echo "Generating coverage HTML report..."
go tool cover -html=coverage/unit.out -o coverage/coverage.html

echo ""
echo "Coverage Summary:"
echo "=================="
go tool cover -func=coverage/unit.out | tail -1

# Extract coverage percentage
COVERAGE=$(go tool cover -func=coverage/unit.out | tail -1 | awk '{print $3}' | sed 's/%//')

echo ""
if (( $(echo "$COVERAGE >= 90" | bc -l) )); then
    echo -e "${GREEN}✓ Coverage target met: ${COVERAGE}% (target: 90%)${NC}"
elif (( $(echo "$COVERAGE >= 80" | bc -l) )); then
    echo -e "${YELLOW}⚠ Coverage acceptable: ${COVERAGE}% (target: 90%)${NC}"
else
    echo -e "${RED}✗ Coverage below target: ${COVERAGE}% (target: 90%)${NC}"
    exit 1
fi

echo ""
echo "Detailed coverage by file:"
echo "=========================="
go tool cover -func=coverage/unit.out | grep -v "total:" | sort -k3 -n

echo ""
echo "Coverage report saved to: coverage/coverage.html"
echo "Open it in a browser to see detailed line-by-line coverage"

# Check for uncovered critical functions
echo ""
echo "Checking for uncovered critical functions..."
CRITICAL_UNCOVERED=$(go tool cover -func=coverage/unit.out | grep -E "(ImportCharacter|Reconcile|Validate|Detect)" | awk '$3 < 80 {print}' || true)

if [ -n "$CRITICAL_UNCOVERED" ]; then
    echo -e "${YELLOW}⚠ Warning: Some critical functions have low coverage:${NC}"
    echo "$CRITICAL_UNCOVERED"
else
    echo -e "${GREEN}✓ All critical functions have adequate coverage${NC}"
fi

echo ""
echo "Running benchmarks..."
go test -bench=. -benchmem ./importer/ -run=^$ | tee coverage/benchmark.txt

echo ""
echo "==================================="
echo "Test Coverage Analysis Complete"
echo "==================================="
