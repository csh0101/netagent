package util

import "fmt"

type Pipeline chan struct{}

type Callback struct {
	name string
	f    func() (chan struct{}, error)
}

func NewCallback(name string, f func() (chan struct{}, error)) Callback {
	return Callback{
		name: name,
		f:    f,
	}
}

func (p Pipeline) Run(callbacks ...Callback) error {
	var name string
	for i, callback := range callbacks {
		if i == 0 {
			close(p)
		}
		// replace it with log
		fmt.Printf("task %s waiting \n", callback.name)
		<-p
		fmt.Printf("task %s starting \n", callback.name)
		f := callback.f
		var err error
		p, err = f()
		if err != nil {
			return err
		}
		name = callback.name
	}
	<-p
	fmt.Println("last task ", name, "has finished!")
	return nil
}
