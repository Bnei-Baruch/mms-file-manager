package main

import (
	"fmt"
	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {
	//	Env = `GOPATH=.vendor::$GOPATH`

	p.Task("default", do.S{"hello", "build"}, nil)

	p.Task("hello", nil, func(c *do.Context) {
		name := c.Args.AsString("name", "n")
		if name == "" {
			c.Bash("echo Hello $USER!")
		} else {
			fmt.Println("Hello", name)
		}
	})


	p.Task("build", do.S{"views", "assets"}, func(c *do.Context) {
		c.Run("GOOS=linux GOARCH=amd64 go build", do.M{"$in": "cmd/server"})
	}).Src("**/*.go")

	p.Task("server", do.S{}, func(c *do.Context) {
		// rebuilds and restarts when a watched file changes
		c.Start("main.go", do.M{"$in": "./"})
	}).Src("**/*.go", "*.{go,json}").
	Debounce(3000)

	p.Task("views", nil, func(c *do.Context) {
		c.Run("razor templates")
	}).Src("templates/**/*.go.html")
}

func main() {
	do.Godo(tasks)
	fmt.Println("11111")
}