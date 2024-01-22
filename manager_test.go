package surveysteps

import (
	"fmt"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.nhat.io/consolesteps"
	"go.nhat.io/surveyexpect"
)

type TestingT struct {
	error *surveyexpect.Buffer
	log   *surveyexpect.Buffer

	clean func()
}

func (t *TestingT) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.error, format, args...)
}

func (t *TestingT) Log(args ...interface{}) {
	_, _ = fmt.Fprintln(t.log, args...)
}

func (t *TestingT) Logf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.log, format, args...)
}

func (t *TestingT) FailNow() {
	panic("failed")
}

func (t *TestingT) Cleanup(clean func()) {
	t.clean = clean
}

func (t *TestingT) ErrorString() string {
	return t.error.String()
}

func (t *TestingT) LogString() string {
	return t.log.String()
}

func T() *TestingT {
	return &TestingT{
		error: new(surveyexpect.Buffer),
		log:   new(surveyexpect.Buffer),
		clean: func() {},
	}
}

func TestManager_ExpectationsWereNotMet(t *testing.T) {
	t.Parallel()

	testingT := T()
	c := consolesteps.New(testingT)
	s := New(testingT).WithConsole(c)
	sc := &godog.Scenario{Id: "42", Name: "ExpectationsWereNotMet"}

	c.NewConsole(sc)

	require.NoError(t, s.expectPasswordAnswer("Enter password:", "password"))

	<-time.After(50 * time.Millisecond)

	s.close(sc)

	expectedError := `in scenario "ExpectationsWereNotMet", there are remaining expectations that were not met:\
[\t\s]*Expect : Password Prompt\
[\t\s]*Message: "Enter password:"\
[\t\s]*Answer : "password"`

	assert.Regexp(t, expectedError, testingT.ErrorString())
}
