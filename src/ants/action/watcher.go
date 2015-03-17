package action

/*
* define a watch
* dead loop for a list
* if there is something,do some action
* implements is in watcher package
* for example
*	*	reporter will loop a result list,if crawler push a result,
*	* 	reporter get the result and send it to master node
**/

type Watcher interface {
	Stop()
	IsStop() bool
	IsPause() bool
	Pause()
	Unpause()
	Start()
	Run()
}
