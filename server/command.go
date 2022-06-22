package server

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	pathSeparator    = string(os.PathSeparator)
	dotPathSeparator = "." + pathSeparator
)

var logCommands []*Command

// populate log commands with default values and args
func init() {
	logCommands = []*Command{
		// ../scep/scepserver-linux-amd64 -allowrenew 0 -challenge nanomdm -debug
		NewCommand().
			WithAlias("scepserver").
			WithName(fmt.Sprintf("scepserver-%s-%s", runtime.GOOS, runtime.GOARCH)).
			//WithName("counter").
			WithPath(fmt.Sprintf(".%sscep%[1]s", pathSeparator)).
			WithArgs("-allowrenew", "0", "-challenge", "nanomdm", "-debug" /*, "-log-json", "2>&1"*/),

		// ngrok http 8080
		// ngrok http 8080 --log=stdout
		NewCommand().
			WithAlias("ngrok_scep").
			WithName("ngrok").
			WithArgs("http", "8080", "--log=stdout"),

		// ./nanomdm-linux-amd64 -ca ca.pem -api nanomdm -debug
		NewCommand().
			WithAlias("nanomdm").
			WithName(fmt.Sprintf("nanomdm-%s-%s", runtime.GOOS, runtime.GOARCH)).
			WithPath(dotPathSeparator).
			WithArgs("-ca", "ca.pem", "-api", "nanomdm", "-debug"),

		// ngrok http 9000
		NewCommand().
			WithAlias("ngrok_nanomdm").
			WithName("ngrok").
			WithArgs("http", "9000"),

		// for tests
		// ping google.com
		NewCommand().
			WithAlias("ping").
			WithName("ping").
			WithArgs("google.com"),
	}
}

type Command struct {
	alias string
	name  string
	path  string
	args  []string
}

// NewCommand creates Command from string if provided as argument
// otherwise creates new empty Command
func NewCommand(s ...string) *Command {
	if len(s) > 0 {
		return NewCommandFromString(strings.Join(s, " "))
	}
	return &Command{}
}

// NewCommandFromString creates Command from string
// if empty creates new empty Command
func NewCommandFromString(s string) *Command {
	cmnd := Command{}
	cmnd.parseString(s)
	return &cmnd
}

func (c *Command) WithAlias(v string) *Command {
	c.alias = v
	return c
}
func (c *Command) WithName(v string) *Command {
	c.name = v
	return c
}
func (c *Command) WithPath(v string) *Command {
	c.path = v
	return c
}
func (c *Command) WithArgs(args ...string) *Command {
	c.args = args
	return c
}

func (c *Command) pathName() string {
	return c.path + c.name
}

// processName first argument in os.StartProcess
func (c *Command) processName() string {
	if c.path == "" {
		return c.name
	}
	return dotPathSeparator + c.name
}

// processDir - os.ProcAttr.Dir in os.StartProcess
func (c *Command) processDir() string {
	if c.path == "" || c.path == dotPathSeparator {
		return ""
	}
	return c.path
}

// pattern use as mux handler pattern
func (c *Command) pattern() string {
	if c.alias != "" {
		return pathSeparator + c.alias
	}
	return pathSeparator + c.name
}

func (c *Command) parseString(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	all := strings.Split(s, " ")
	c.path, c.name = path.Split(all[0])
	c.args = all[1:]
}

func (c *Command) String() string {
	return c.path + c.name + " " + strings.Join(c.args, " ")
}
