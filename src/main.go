package main

/* Purpose: lightly wrap the os.CMD to preform long running processes that can be canceled at any time.
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
	/*	//args := []string{}
		args := []string{"-c", `for((i=1;i<=10;i+=2)); do echo "Welcome $i times"; sleep 1; done`}
		proc := process.NewProcess("bash", args...)
		result := make(chan error)

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

		time.Sleep(time.Second * 3)
		proc.Kill()
		select {
		case retCode := <-result:
			log.Println(retCode)
		}
	*/
}
