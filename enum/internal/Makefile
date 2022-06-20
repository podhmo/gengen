default: gen.go

gen.go: gen.go.tmpl main.go sandbox/*.go
	go run main.go --packages github.com/podhmo/gengen/enum/internal/sandbox --interface github.com/podhmo/gengen/enum.Interface > $@