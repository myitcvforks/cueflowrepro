package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

// TestMain allows us to setup "main" as a special
// command within the testscripts
func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"main": main1,
	}))
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:           "testdata",
		UpdateScripts: os.Getenv("TESTSCRIPT_UPDATE") != "",
	})

}
