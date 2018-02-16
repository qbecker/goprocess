package process

import (
	//"log"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	args := []string{"Hello, World!"}
	proc := NewProcess("echo", args...)
	answer := "Hello, World!"
	result := ""

	out := proc.StreamOutput()
	proc.Start()
	go func() {
		for out.Scan() {
			result = result + out.Text()
		}
		if result != answer {

			t.Errorf("Incorrect output, expected: %s, got: %s", answer, result)
		}
	}()
}

func TestCancel(t *testing.T) {
	args := []string{"while true; do foo; done"}
	proc := NewProcess("bash", args...)
	result := make(chan error)

	go func() {
		result <- proc.Wait()
	}()
	proc.Start()

	time.Sleep(time.Second * 1)
	proc.Kill()
	select {
	case retCode := <-result:
		if retCode == nil {
			t.Errorf("Incorrect")
		}
	}
}
func TestDoubleKill(t *testing.T) {
	args := []string{"while true; do foo; done"}
	proc := NewProcess("bash", args...)
	result := make(chan error)

	go func() {
		result <- proc.Wait()
	}()
	proc.Start()

	time.Sleep(time.Second * 1)
	proc.Kill()
	select {
	case retCode := <-result:
		if retCode == nil {
			t.Errorf("Incorrect")
		}
	}
	proc.Kill()
}

func TestStreamAfterStart(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	args := []string{"Hello, World!"}
	proc := NewProcess("echo", args...)
	proc.Start()
	proc.StreamOutput()
}
func TestDoubleInStream(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	args := []string{"Hello, World!"}
	proc := NewProcess("echo", args...)

	proc.StreamOutput()
	proc.StreamOutput()
	proc.Start()
}