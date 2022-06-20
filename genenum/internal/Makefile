default: gen.go

gen.go: gen.go.tmpl main.go
	go run main.go --packages github.com/podhmo/gengen/genenum/sandbox --interface github.com/podhmo/gengen/genenum/schema.EnumInterface > $@