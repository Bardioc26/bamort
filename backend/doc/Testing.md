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