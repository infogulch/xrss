package main

import (
	"html/template"

	"github.com/infogulch/xtemplate"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	funcmap := template.FuncMap{"fetchFeed": FetchFeed}
	xtemplate.Main(xtemplate.WithFuncMaps(funcmap))
}
