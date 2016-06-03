package lib

import (
	"fmt"
	"bufio"
	"bytes"
)

func (s *Seekret) Inspect() {
	for _, o := range s.objectList {
		if len(o.Content) < MaxObjectContent {
			for _, r := range s.ruleList {
				x := bufio.NewScanner(bytes.NewReader(o.Content))
				buf := []byte{}
				x.Buffer(buf, MaxObjectContent)

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
						if unmatch == false {
							secret := Secret{
								Object:    o,
								Rule:      r,
								Nline:     nLine,
								Line:      line,
							}
							secret.Exception = exceptionCheck(s.exceptionList, secret)
							s.secretList = append(s.secretList, secret)

						}
					}
				}

				if err := x.Err(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}