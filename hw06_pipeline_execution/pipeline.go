package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	worker := func(done In, chIn In) Out {
		ch := make(chan interface{})
		go func() {
			defer close(ch)
			for {
				select {
				case <-done:
					return
				case v, ok := <-chIn:
					if !ok {
						return
					}
					ch <- v
				}
			}
		}()
		return ch
	}

	for _, s := range stages {
		out = worker(done, s(out))
	}

	return out
}
