package seekret

import (
	"fmt"
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


func inspect_worker(id int, jobs <-chan workerJob, results chan<- workerResult) {
	lid := 0
	for job := range jobs {
		lid = lid +1
		result := workerResult{
			wid: id,
		}

		content := job.objectGroup[0].Content

		for _,r := range job.ruleList {
			if r.Enabled == false {
				continue
			}

			fs := bufio.NewScanner(bytes.NewReader(content))
			buf := []byte{}

			// INFO: Remove the next two lines if using golang < 1.6
			fs.Buffer(buf, models.MaxObjectContentLen)

			nLine := 0
			for fs.Scan() {
				nLine = nLine + 1
				line := fs.Text()

				runResultList := r.Run([]byte(line))

				for _,object := range job.objectGroup {
					for _,runResult := range runResultList {
						secret := models.NewSecret(&object, &r, runResult.Nline, runResult.Line)
						secret.SetException(exceptionCheck(job.exceptionList, *secret))
						result.secretList = append(result.secretList, *secret)
					}
				}
			}

		}

		results <- result
	}
}

func (s *Seekret)Inspect(workers int) {
	jobs := make(chan workerJob)
	results := make(chan workerResult)

	for w := 1; w <= workers; w++ {
		go inspect_worker(w, jobs, results)
	}

	objectGroupMap := s.GroupObjectsByPrimaryKeyHash()
	go func() {
		for _,objectGroup := range objectGroupMap {
			jobs <- workerJob{
				objectGroup: objectGroup,
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

