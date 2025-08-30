# Run go Tests

## Benchmarks

### Run all benchmarks
go test -bench=. -benchmem ./models

### Run specific category
go test -bench=BenchmarkSource -benchmem ./models

### Run lightweight benchmarks only
go test -bench=BenchmarkSimple -benchmem ./models

### Run parallel benchmarks
go test -bench=Parallel -benchmem ./models

## Coverage

go test ./equipment/... -cover
go test ./equipment/... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
### eval Coverage
    mode: set
    bamort/equipment/handlers.go:18.67,20.2 1 1
    bamort/equipment/handlers.go:22.40,24.55 2 1
    bamort/equipment/handlers.go:24.55,27.3 2 1
    bamort/equipment/handlers.go:29.2,29.63 1 1
    bamort/equipment/handlers.go:29.63,32.3 2 0
    bamort/equipment/handlers.go:34.2,34.41 1 1
    bamort/equipment/handlers.go:37.38,41.100 3 1
    bamort/equipment/handlers.go:41.100,44.3 2 0
    bamort/equipment/handlers.go:46.2,46.36 1 1
    bamort/equipment/handlers.go:49.40,53.77 3 1
    bamort/equipment/handlers.go:53.77,56.3 2 1
    bamort/equipment/handlers.go:58.2,58.55 1 1
    bamort/equipment/handlers.go:58.55,61.3 2 1
    bamort/equipment/handlers.go:63.2,63.61 1 1
    bamort/equipment/handlers.go:63.61,66.3 2 0
    bamort/equipment/handlers.go:68.2,68.36 1 1
    bamort/equipment/handlers.go:71.40,73.89 2 1
    bamort/equipment/handlers.go:73.89,76.3 2 1
    bamort/equipment/handlers.go:78.2,78.77 1 1

Looking at the coverage.out file, you can see which lines are covered. The format is:
file:startLine.startCol,endLine.endCol numStatements covered

Where the last number indicates if the line is covered:
1 = covered
0 = not covered

From your coverage.out, I can see these uncovered lines (marked with 0):
Line 29-32: bamort/equipment/handlers.go:29.63,32.3 2 0 - Database error handling in CreateAusruestung
Line 41-44: bamort/equipment/handlers.go:41.100,44.3 2 0 - Database error handling in ListAusruestung
Line 63-66: bamort/equipment/handlers.go:63.61,66.3 2 0 - Database save error handling in UpdateAusruestung 