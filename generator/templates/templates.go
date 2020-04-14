package templates

const GitignoreTpl string = `# Created by https://www.gitignore.io/api/go,intellij+iml

### Go ###
# Compiled Object files, Static and Dynamic libs (Shared Objects)
*.o
*.a
*.so

# Folders
_obj
_test

# Architecture specific extensions/prefixes
*.[568vq]
[568vq].out

*.cgo1.go
*.cgo2.c
_cgo_defun.c
_cgo_gotypes.go
_cgo_export.*

_testmain.go

*.exe
*.test
*.prof

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# external packages folder
vendor/

# IntelliJ
*.iml
*.ipr
*.iws
modules.xml
.idea
.idea_modules
#/out/

## Plugin-specific files:

# JIRA plugin
atlassian-ide-plugin.xml

# Crashlytics plugin (for Android Studio and IntelliJ)
com_crashlytics_export_strings.xml
crashlytics.properties
crashlytics-build.properties
fabric.properties

# Sublime text
*.sublime-*

# Coverage
.cover
`
const MainTpl string = `package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"%s"
)

func main() {

	app := cli.NewApp()
	app.Name = "%s"
	app.Usage = "%s"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{}

	app.Action = func(c *cli.Context) error {
    	fmt.Println("hello, world")
    	return nil
  	}

	app.Run(os.Args)
}
`
