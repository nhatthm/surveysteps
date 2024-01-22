package surveysteps

import (
	"errors"
	"sync"

	"github.com/stretchr/testify/require"
	"go.nhat.io/surveyexpect"
)

// Survey is a wrapper around *surveyexpect.Survey to make it run with cucumber/godog.
type Survey struct {
	*surveyexpect.Survey

	test surveyexpect.TestingT
	mu   sync.Mutex

	doneChan chan struct{}
}

func (s *Survey) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.getDoneChanLocked()
}

func (s *Survey) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}

	return s.doneChan
}

func (s *Survey) closeDoneChan() {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := s.getDoneChanLocked()

	select {
	case <-ch:
		// Already closed. Don't close again.

	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(ch)
	}
}

// Expect runs an expectation against a given console.
func (s *Survey) Expect(c surveyexpect.Console) error {
	for {
		select {
		case <-s.getDoneChan():
			return nil

		default:
			err := s.Survey.Expect(c)
			if err != nil && !errors.Is(err, surveyexpect.ErrNothingToDo) {
				return err
			}
		}
	}
}

// Start starts a new survey.
func (s *Survey) Start(console surveyexpect.Console) *Survey {
	go func() {
		require.NoError(s.test, s.Expect(console))
	}()

	return s
}

// Close notifies other parties and close the survey.
func (s *Survey) Close() {
	s.closeDoneChan()
}

// NewSurvey creates a new survey.
func NewSurvey(t surveyexpect.TestingT, options ...surveyexpect.ExpectOption) *Survey {
	return &Survey{
		Survey: surveyexpect.New(t, options...),
		test:   t,
	}
}
