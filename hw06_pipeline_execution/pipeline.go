package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// check for empty in channel or empty stages slice
	if in == nil || len(stages) == 0 {
		return in
	}

	worker := func(in, done In) Out {
		out := make(Bi)
		go func() {
			// closes out channel and reads remaining in channel data
			defer func() {
				close(out)
				for range in {
					<-in
				}
			}()
			for {
				select {
				case <-done: // check for interrupt signal
					return
				case value, ok := <-in: // read values from in channel until channel is empty
					if !ok {
						return
					}
					out <- value
				}
			}
		}()
		return out
	}

	// init stages
	out := in
	for _, stage := range stages {
		out = worker(stage(out), done)
	}

	return out
}
