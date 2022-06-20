default: gen.go

gen.go: gen.go.tmpl main.go sandbox/*.go
	go run main.go --packages github.com/podhmo/gengen/genenum/internal/sandbox --interface github.com/podhmo/gengen/genenum/schema.Interface > $@