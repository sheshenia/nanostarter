package server

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
)

var logCommands []*Command

func init() {
	const (
		osN   = runtime.GOOS
		archN = runtime.GOARCH
	)

	//populate log commands with default values and args
	logCommands = []*Command{
		// ./scepserver-linux-amd64 -allowrenew 0 -challenge nanomdm -debug
		NewCommand().
			WithAlias("scepserver").
			WithName(fmt.Sprintf("scepserver-%s-%s", runtime.GOOS, runtime.GOARCH)).
			WithPath("."+string(os.PathSeparator)).
			WithArgs("-allowrenew", "0", "-challenge", "nanomdm", "-debug"),

		// ngrok http 8080
		NewCommand().
			WithAlias("ngrok_scep").
			WithName("ngrok").
			WithArgs("http", "8080"),

		// ./nanomdm-linux-amd64  -ca ca.pem -api nanomdm -debug
		NewCommand().
			WithAlias("nanomdm").
			WithName(fmt.Sprintf("nanomdm-%s-%s", runtime.GOOS, runtime.GOARCH)).
			WithPath("."+string(os.PathSeparator)).
			WithArgs("-ca", "ca.pem", "-api", "nanomdm", "-debug"),

		// ngrok http 9000
		NewCommand().
			WithAlias("ngrok_nanomdm").
			WithName("ngrok").
			WithArgs("http", "9000"),
	}
}

type Command struct {
	alias string
	name  string
	path  string
	args  []string
}

func (c *Command) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.Write([]byte(c.alias))
}

// NewCommand creates Command from string if provided as argument
// otherwise creates new empty Command
func NewCommand(s ...string) *Command {
	if len(s) > 0 {
		return NewCommandFromString(strings.Join(s, " "))
	}
	return &Command{}
}

func NewCommandFromString(s string) *Command {
	s = strings.TrimSpace(s)
	cmnd := Command{}

	// if empty string returns empty Command
	if s == "" {
		return &cmnd
	}

	all := strings.Split(s, " ")
	cmnd.path, cmnd.name = path.Split(all[0])
	cmnd.alias = cmnd.name //by default alias equal to name
	cmnd.args = all[1:]
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
