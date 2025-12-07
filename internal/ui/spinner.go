package ui

import (
	"fmt"
	"sync"
	"time"
)

// Spinner shows an animated loading indicator
type Spinner struct {
	message  string
	active   bool
	done     chan bool
	mu       sync.Mutex
	frames   []string
	interval time.Duration
}

// NewSpinner creates a new loading spinner with a message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:  message,
		active:   false,
		done:     make(chan bool),
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		interval: 100 * time.Millisecond,
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go s.animate()
}

// Stop halts the spinner and shows completion message
func (s *Spinner) Stop(success bool) {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	s.done <- true
	<-s.done // Wait for animation to finish

	// Clear the line and show final status
	fmt.Print("\r")
	if success {
		PrintSuccess(s.message)
	} else {
		PrintError(s.message)
	}
}

// StopWithMessage stops the spinner and shows a custom message
func (s *Spinner) StopWithMessage(message string, success bool) {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	s.done <- true
	<-s.done // Wait for animation to finish

	// Clear the line and show final status
	fmt.Print("\r")
	if success {
		PrintSuccess(message)
	} else {
		PrintError(message)
	}
}

// animate runs the spinner animation loop
func (s *Spinner) animate() {
	frameIndex := 0
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			// Clear the spinner line
			fmt.Print("\r" + clearLine())
			s.done <- true
			return
		case <-ticker.C:
			s.mu.Lock()
			if !s.active {
				s.mu.Unlock()
				return
			}
			s.mu.Unlock()

			frame := s.frames[frameIndex]
			frameIndex = (frameIndex + 1) % len(s.frames)

			// Print spinner with message
			fmt.Printf("\r%s %s %s",
				ColorCyan.Sprint(frame),
				s.message,
				clearLine(),
			)
		}
	}
}

// clearLine returns ANSI escape code to clear rest of line
func clearLine() string {
	return "\033[K"
}

// ShowSpinner is a convenience function to show a spinner for a function execution
func ShowSpinner(message string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()

	err := fn()

	if err != nil {
		spinner.Stop(false)
	} else {
		spinner.Stop(true)
	}

	return err
}

// ShowSpinnerWithResult shows spinner and returns both result and error
func ShowSpinnerWithResult(message string, fn func() (interface{}, error)) (interface{}, error) {
	spinner := NewSpinner(message)
	spinner.Start()

	result, err := fn()

	if err != nil {
		spinner.Stop(false)
	} else {
		spinner.Stop(true)
	}

	return result, err
}
