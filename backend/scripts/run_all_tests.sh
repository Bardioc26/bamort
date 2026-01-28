set -e
cd /data/dev/bamort/backend
echo "dont forget to reactivate skipped tests after fixing issues"
go test ./cmd -v |grep FAIL
go test ./database -v |grep FAIL
go test ./maintenance -v |grep FAIL
go test ./testutils -v |grep FAIL
go test ./router -v |grep FAIL
go test ./models -v |grep FAIL
go test ./api -v |grep FAIL
go test ./gamesystem -v |grep FAIL
go test ./transfer -v |grep FAIL
go test ./uploads -v |grep FAIL
go test ./user -v |grep FAIL
go test ./importer -v |grep FAIL
go test ./character -v |grep FAIL
go test ./gsmaster -v |grep FAIL
go test ./pdfrender -v |grep FAIL
go test ./config -v |grep FAIL
go test ./equipment -v |grep FAIL
go test ./logger    -v |grep FAIL

# Optional: generate coverage report for the whole backend module.
# Enable by setting RUN_COVERAGE=1 in the environment.
if [ "${RUN_COVERAGE:-}" = "1" ]; then
	echo "Running coverage for backend (this may take a while)..."
	# produce a single combined coverage profile
	go test ./... -coverprofile=coverage.out -covermode=atomic
	if [ -f coverage.out ]; then
		# generate an HTML report (best viewed in a browser)
		go tool cover -html=coverage.out -o coverage.html || true
		echo "Coverage written to coverage.out and coverage.html"
	else
		echo "coverage.out not created"
	fi
fi