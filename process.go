package process

/* Purpose: lightly wrap the os.CMD to preform long running processes that can be canceled at any time.
   Requirments:
       Processes must be able to be canceled at any time.
       Standard out must be captured.
       Standard Error must be captured.
   S-Goal(s):
           Ability to pipe in standard in at any time.
*/
type Process struct {
}
