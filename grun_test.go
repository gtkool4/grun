package grun_test

import (
	"errors"
	"testing"

	"github.com/gtkool4/grun"
)

func Test_errorPaths(t *testing.T) {
	var errs grun.Errors
	if errs.IsError() || errs.Error() != "" {
		t.Fail()
	}

	errs.Append(errors.New("fail"))
	if errs.ToError().Error() != "fail" {
		t.Fail()
	}
}

func Test_newApps(t *testing.T) {
	list := map[string]*grun.App{
		"Tiny":   grun.NewTiny(),
		"Small":  grun.NewSmall(),
		"Medium": grun.NewMedium(),
		"Large":  grun.NewLarge()}
	for name, appInfo := range list {
		if appInfo == nil || appInfo.Width < 1 || appInfo.Height < 1 {
			t.Errorf("generator failed for size: %s\n", name)
		}
	}
}
