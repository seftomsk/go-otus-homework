package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runStage(in In, done In) Out {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case stream <- val:
				}
			}
		}
	}()
	return stream
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(runStage(in, done))
	}
	return in
}
