package main

import (
	"github.com/infogulch/xtemplate"

	_ "github.com/infogulch/xtemplate-gofeed"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	xtemplate.Main()
}
