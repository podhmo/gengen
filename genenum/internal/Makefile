default: gen.go

gen.go: gen.go.tmpl main.go op.go
	go run main.go --packages github.com/podhmo/gengen/genenum/internal --interface github.com/podhmo/gengen/genenum/schema.Interface > $@