package bootstrap

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.nhat.io/surveyexpect"

	"go.nhat.io/surveysteps"
)

// Used by init().
//
//nolint:gochecknoglobals
var (
	runGoDogTests bool

	out = new(surveyexpect.Buffer)
	opt = godog.Options{
		Strict: true,
		Output: out,
	}
)

// This has to run on init to define -godog flag, otherwise "undefined flag" error happens.
//
//nolint:gochecknoinits
func init() {
	flag.BoolVar(&runGoDogTests, "godog", false, "Set this flag is you want to run godog BDD tests")
	godog.BindCommandLineFlags("", &opt)
}

func TestIntegration(t *testing.T) {
	if !runGoDogTests {
		t.Skip(`Missing "-godog" flag, skipping integration test.`)
	}

	p := NewPrompt()
	m := surveysteps.New(t).
		WithStarter(p.WithStdio)

	RunSuite(t, "..", func(_ *testing.T, ctx *godog.ScenarioContext) {
		m.RegisterSteps(ctx)
		p.RegisterSteps(ctx)
	})
}

func RunSuite(t *testing.T, path string, featureContext func(t *testing.T, ctx *godog.ScenarioContext)) {
	t.Helper()

	flag.Parse()

	if opt.Randomize == 0 {
		opt.Randomize = rand.Int63() // nolint: gosec
	}

	var paths []string

	files, err := os.ReadDir(filepath.Clean(path))
	require.NoError(t, err)

	paths = make([]string, 0, len(files))

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".feature") {
			paths = append(paths, filepath.Join(path, f.Name()))
		}
	}

	for _, path := range paths {
		path := path

		t.Run(path, func(t *testing.T) {
			opt.Paths = []string{path}
			suite := godog.TestSuite{
				Name:                 "Integration",
				TestSuiteInitializer: nil,
				ScenarioInitializer: func(s *godog.ScenarioContext) {
					featureContext(t, s)
				},
				Options: &opt,
			}
			status := suite.Run()

			if status != 0 {
				fmt.Println(out.String())
				assert.Fail(t, "one or more scenarios failed in feature: "+path)
			}
		})
	}
}
