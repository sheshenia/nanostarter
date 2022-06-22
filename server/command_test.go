package server

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestNewCommand(t *testing.T) {
	tCmd := NewCommand().
		WithAlias("nanomdm").
		WithName(fmt.Sprintf("nanomdm-%s-%s", runtime.GOOS, runtime.GOARCH)).
		WithPath(dotPathSeparator).
		WithArgs("-ca", "ca.pem", "-api", "nanomdm", "-debug")

	exp := fmt.Sprintf(`.%snanomdm-%s-%s -ca ca.pem -api nanomdm -debug`,
		pathSeparator, runtime.GOOS, runtime.GOARCH)

	t.Run("test cmd.String()", func(t *testing.T) {
		if tCmd.String() != exp {
			t.Errorf("exp: %s\ngot: %s\n", exp, tCmd.String())
		}
	})

	t.Run("test cmd.pathName()", func(t *testing.T) {
		expPathName := dotPathSeparator + fmt.Sprintf("nanomdm-%s-%s", runtime.GOOS, runtime.GOARCH)
		if tCmd.pathName() != expPathName {
			t.Errorf("exp: %s\ngot: %s\n", expPathName, tCmd.pathName())
		}
	})

	// 1) "./" => ""
	// 2) "" => ""
	// 3) "./test/ => "./test/"
	// initial path = dotPathSeparator

	// 1)
	t.Run("test cmd.processDir() 1", func(t *testing.T) {
		if tCmd.processDir() != "" {
			t.Errorf("exp: %s\ngot: %s\n", "", tCmd.processDir())
		}
	})
	// 2)
	t.Run("test cmd.processDir() 2", func(t *testing.T) {
		tCmd.path = ""
		if tCmd.processDir() != "" {
			t.Errorf("exp: %s\ngot: %s\n", "", tCmd.processDir())
		}
	})
	// 3)
	t.Run("test cmd.processDir() 3", func(t *testing.T) {
		expPD := dotPathSeparator + "test" + pathSeparator
		tCmd.path = expPD
		if tCmd.processDir() != expPD {
			t.Errorf("exp: %s\ngot: %s\n", expPD, tCmd.processDir())
		}
	})
	// restore initial path
	tCmd.WithPath(dotPathSeparator)

	// now vice-versa, try to parse from string, and to compare with tCmd we have
	tCmd2 := NewCommand(exp) //with arg it will call NewCommandFromString
	t.Run("test NewCommandFromString", func(t *testing.T) {
		if reflect.DeepEqual(*tCmd, *tCmd2) {
			t.Error("not deep equal tCmd and tCmd2")
		}
	})

	t.Run("test cmd.pattern() with alias", func(t *testing.T) {
		exp := pathSeparator + tCmd.alias
		if tCmd.pattern() != exp {
			t.Errorf("exp: %s\ngot: %s\n", exp, tCmd.pattern())
		}
	})

	t.Run("test cmd.pattern() without alias", func(t *testing.T) {
		tCmd.alias = ""
		exp := pathSeparator + tCmd.name
		if tCmd.pattern() != exp {
			t.Errorf("exp: %s\ngot: %s\n", exp, tCmd.pattern())
		}
	})
}
