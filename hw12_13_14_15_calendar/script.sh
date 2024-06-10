goose -dir ./migrations postgres "postgres://postgres:1234512345@localhost:5432/calendardb?sslmode=disable" up
#goose -dir ./migrations postgres "postgres://postgres:1234512345@localhost:5432/calendardb?sslmode=disable" down
#find . -name '*.go' | xargs gci write --skip-generated -s standard -s default
#find . -name '*.go' | xargs gofumpt -l -w