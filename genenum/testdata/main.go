//go:build gen
// +build gen

package main

type Schema struct {
	PackageName string
}

tmpl := `package {{.PackageName}}
`


func main() {
	template.New	
}
