package main

import (
	"fmt"
	"github.com/redneckbeard/quimby"
	"github.com/redneckbeard/gadget/strutil"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	quimby.Add(&Controller{})
}

type Controller struct {
	*quimby.Flagger
	name string
}

func (c *Controller) Desc() string {
	return "Create a new controller. Camel-cased names will be used unmodified. Other names will be title-cased."
}

func (c *Controller) SetFlags() {
	c.StringVar(&c.name, "name", "", "Name of the controller type to create")
}

func (c *Controller) Run() {
	name := c.name
	if name == "" {
		args := c.FlagSet.Args()
		if len(args) == 0 {
			fmt.Println("You must supply a name for the controller.")
			return
		}
		name = args[0]
	}
	current, _ := os.Getwd()
	path, err := filepath.Rel(
		filepath.Join(os.ExpandEnv("$GOPATH"), "src"),
		filepath.Join(current),
	)
	if err != nil {
		fmt.Println("Projects must be created in the src directory of your GOPATH.")
	}
	createControllerFile(path, name)
}

func createControllerFile(projectName, controllerName string) {
	var filename string
	t := getTemplate("controller.tpl")
	if strutil.IsPascalCase(controllerName) {
		filename = strutil.Snakify(controllerName)
	} else {
		filename = strings.ToLower(controllerName)
		controllerName = strings.Title(controllerName)
	}
	path := filepath.Join(projectName, "controllers", filename + ".go")
	if f, err := os.Create(path); err != nil {
		fmt.Printf("Unable to create file controllers/%s.go: %s\n", filename, err)
	} else {
		defer f.Close()
		t.Execute(f, map[string]string{
			"name": controllerName,
			"project": projectName,
		})
		fmt.Printf("Created %s\n", path)
	}
}
