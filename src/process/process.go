package process

/* Purpose: lightly wrap the os.CMD to preform long running processes that can be canceled at any time(and hopfully
    be given input at any time.
   Requirments:
       Processes must be able to be canceled at any time.
       Standard out must be captured.
       Standard Error must be captured.
   S-Goal(s):
           Ability to pipe in stdin at any time.
*/

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
)

type Process struct {
	proc               *exec.Cmd
	cancellationSignal chan uint8
	done               chan error
	returnCode         chan error
	started            bool
	stdOutRead         *io.PipeReader
	stdOutWrite        *io.PipeWriter
	inputWriter        *io.PipeWriter
	inputStreamSet     bool
	outputStreamSet    bool
	// Access to completed MUST capture the lock.
	completed bool
	mutex     sync.RWMutex
}

func NewProcess(command string, args ...string) *Process {
	process := &Process{
		exec.Command(command, args...),
		make(chan uint8, 1),
		make(chan error, 1),
		make(chan error, 1),
		false,
		&io.PipeReader{},
		&io.PipeWriter{},
		&io.PipeWriter{},
		false,
		false,
		false,
		sync.RWMutex{}}
	return process
}

func (p *Process) Start() *Process {
	p.started = true
	//Call the other functions to stream stdin and stdout
	err := p.proc.Start()
	if err != nil {
		panic(err)
	}
	go p.awaitOutput()
	go p.finishTimeOutOrDie()
	return p
}

func (p *Process) Wait() error {
	return <-p.returnCode
}

func (p *Process) awaitOutput() {
	//send the exit code to the done channel to signify the command finished
	p.done <- p.proc.Wait()
}

func (p *Process) Kill() {
	p.mutex.Lock()
	if !p.completed {
		p.cancellationSignal <- 1
	}
	p.mutex.Unlock()
}

func (p *Process) OpenInputStream() (io.WriteCloser, error) {
	if p.inputStreamSet {
		panic("Input stream already set")
	}
	if p.started {
		panic("process already started")
	}
	stdIn, err := p.proc.StdinPipe()
	p.inputStreamSet = true
	return stdIn, err

}
func (p *Process) StreamOutput() *bufio.Scanner {
	//pipe both stdout and stderr into the same pipe
	//panics if you do streamoutput after proccess starting or
	//if the output stream is already set
	if p.started {
		panic("Cant set output stream after starting")
	}
	if p.outputStreamSet {
		panic("output stream already set!")
	}
	p.stdOutRead, p.stdOutWrite = io.Pipe()
	p.proc.Stdout = p.stdOutWrite
	p.proc.Stderr = p.stdOutWrite
	p.outputStreamSet = true
	//return a scanner which they can read from till empty
	return bufio.NewScanner(p.stdOutRead)
}

func (p *Process) finishTimeOutOrDie() {
	defer p.cleanup()
	var result error
	select {
	case result = <-p.done:
	case <-p.cancellationSignal:
		log.Println("received cancellationSignal")
		//NOT PORTABLE TO WINDOWS
		err := p.proc.Process.Kill()
		if err != nil {
			log.Println(err)
		}
	}
	p.returnCode <- result
}

func (p *Process) cleanup() {
	p.mutex.Lock()
	p.completed = true
	p.mutex.Unlock()
	if p.outputStreamSet {
		p.stdOutRead.Close()
		p.stdOutWrite.Close()
	}
	close(p.done)
	close(p.cancellationSignal)
}
