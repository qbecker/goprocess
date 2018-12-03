package main

/* Purpose: lightly wrap the os.CMD to perform long running processes that can be canceled at any time.
   Requirments:
       Processes must be able to be canceled at any time.
       Standard out must be captured.
       Standard Error must be captured.
   S-Goal(s):
           Ability to pipe in stdin at any time.
*/

import (
	"./process"
	"log"
	"time"
)

func main() {
	//args := []string{}
	args := []string{"if=/dev/zero", "of=test", "bs=1000M", "count=1"}
	proc := process.NewProcess("dd", args...)
	result := make(chan error)
	proc.SetTimeout(time.Second * 3)
	go func() {
		result <- proc.Wait()
	}()

	out := proc.StreamOutput()
	go func() {
		for out.Scan() {
			log.Println(out.Text())
		}
	}()
	proc.Start()
	select {
	case retCode := <-result:
		log.Println(retCode)
	}

}
