// Copyright 2016 - Authors included on AUTHORS file.
//
// Use of this source code is governed by a Apache License
// that can be found in the LICENSE file.

package seekret

import (
	"bufio"
	"bytes"
	"github.com/apuigsech/seekret/models"
)

type workerJob struct {
	objectGroup   []models.Object
	ruleList      []models.Rule
	exceptionList []models.Exception
}

type workerResult struct {
	wid        int
	secretList []models.Secret
}

// Inspect executes the inspection into all loaded objects, by checking all
// rules and exceptions loaded.
func (s *Seekret) Inspect(Nworkers int) {
	jobs := make(chan workerJob)
	results := make(chan workerResult)

	for w := 1; w <= Nworkers; w++ {
		go inspect_worker(w, jobs, results)
	}

	objectGroupMap := s.GroupObjectsByPrimaryKeyHash()

	go func() {
		for _, objectGroup := range objectGroupMap {
			jobs <- workerJob{
				objectGroup:   objectGroup,
				ruleList:      s.ruleList,
				exceptionList: s.exceptionList,
			}
		}
		close(jobs)
	}()

	for i := 0; i < len(objectGroupMap); i++ {
		result := <-results
		s.secretList = append(s.secretList, result.secretList...)
	}
}

func inspect_worker(id int, jobs <-chan workerJob, results chan<- workerResult) {
	lid := 0
	for job := range jobs {
		lid = lid + 1
		result := workerResult{
			wid: id,
		}

		content := job.objectGroup[0].Content

		for ri,r := range job.ruleList {
			if r.Enabled == false {
				continue
			}

			fs := bufio.NewScanner(bytes.NewReader(content))
			buf := []byte{}

			// INFO: Remove the next line if using golang < 1.6
			fs.Buffer(buf, models.MaxObjectContentLen)

			nLine := 0
			for fs.Scan() {
				nLine = nLine + 1
				line := fs.Text()

				runResult := r.Run(line)

				for oi,_ := range job.objectGroup {
					if runResult {
						secret := models.NewSecret(&job.objectGroup[oi], &job.ruleList[ri], nLine, line)
						secret.SetException(exceptionCheck(job.exceptionList, *secret))
						result.secretList = append(result.secretList, *secret)
					}
				}

			}
		}

		results <- result
	}
}

func exceptionCheck(exceptionList []models.Exception, secret models.Secret) bool {
	for _, x := range exceptionList {
		match := x.Run(&secret)

		if match {
			return true
		}
	}
	return false
}
