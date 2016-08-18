package seekret

import (
	"bufio"
	"bytes"
	"github.com/apuigsech/seekret/models"
)

type workerJob struct {
	object        models.Object
	ruleList      []models.Rule
	exceptionList []Exception
}

type workerResult struct {
	wid        int
	secretList []Secret
}

func inspect_worker(id int, jobs <-chan workerJob, results chan<- workerResult) {
	for job := range jobs {
		result := workerResult{
			wid: id,
		}

		for _, r := range job.ruleList {
			if r.Enabled == false {
				continue
			}
			x := bufio.NewScanner(bytes.NewReader(job.object.Content))
			buf := []byte{}

			// INFO: Remove the next two lines if using golang < 1.6
			x.Buffer(buf, models.MaxObjectContentLen)

			nLine := 0
			for x.Scan() {
				nLine = nLine + 1
				line := x.Text()

				if r.Match.MatchString(line) {
					unmatch := false
					for _, Unmatch := range r.Unmatch {
						if Unmatch.MatchString(line) {
							unmatch = true
						}
					}
					if !unmatch {
						secret := Secret{
							Object: job.object,
							Rule:   r,
							Nline:  nLine,
							Line:   line,
						}
						secret.Exception = exceptionCheck(job.exceptionList, secret)
						result.secretList = append(result.secretList, secret)
					}
				}
			}
		}
		results <- result
	}
}

func (s *Seekret) Inspect(workers int) {
	jobs := make(chan workerJob)
	results := make(chan workerResult)

	for w := 1; w <= workers; w++ {
		go inspect_worker(w, jobs, results)
	}

	go func() {
		for _, o := range s.objectList {
			jobs <- workerJob{
				object:        o,
				ruleList:      s.ruleList,
				exceptionList: s.exceptionList,
			}
		}
		close(jobs)
	}()

	for i := 0; i < len(s.objectList); i++ {
		result := <-results
		s.secretList = append(s.secretList, result.secretList...)
	}
}
