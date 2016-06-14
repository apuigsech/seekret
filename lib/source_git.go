package lib

import (
	"fmt"
	"github.com/jasonmoo/ssh_config"
	"github.com/libgit2/git2go"
	"io/ioutil"
	"path/filepath"
	"net/url"
	"os"
	"regexp"
)

var (
	SourceTypeGit = &SourceGit{}
)

type SourceGit struct{}

type SourceGitLoadOptions struct {
	Count int
}

func prepareGitLoadOptions(o LoadOptions) SourceGitLoadOptions {
	opt := SourceGitLoadOptions{
		Count: 0,
	}

	if count, ok := o["count"].(int); ok == true {
		opt.Count = count
	}

	return opt
}

func (s *SourceGit) LoadObjects(source string, o LoadOptions) ([]Object, error) {
	var objectList []Object

	opt := prepareGitLoadOptions(o)

	repo, err := openGitRepo(source)
	if err != nil {
		return nil, err
	}

	walk, err := repo.Walk()
	if err != nil {
		return nil, err
	}

	if opt.Count > 0 {
		walk.PushRange(fmt.Sprintf("HEAD~%d..HEAD", opt.Count))
	} else {
		walk.PushHead()
	}
	walk.Sorting(git.SortTime)

	err = walk.Iterate(func(commit *git.Commit) bool {
		tree, err := commit.Tree()
		if err != nil {
			fmt.Println(err)
		}

		// TODO: what to return?
		tree.Walk(func(base string, tentry *git.TreeEntry) int {
			if tentry.Type == git.ObjectBlob {
				blob, err := repo.LookupBlob(tentry.Id)
				if err != nil {
					return 0
				}
				o := Object{
					Name:    fmt.Sprintf("%s:%s%s[%s]", commit.Id().String(), base, tentry.Name, tentry.Id.String()),
					Content: blob.Contents(),
				}
				objectList = append(objectList, o)
			}

			return 0
		})

		return true
	})

	if err != nil {
		return nil, err
	}

	return objectList, nil
}

func credentialsCallback(gitUri string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	sshConfigFile := os.ExpandEnv("$HOME/.ssh/config")

	fh, err := os.Open(sshConfigFile)
	if err != nil {
		panic(err)
	}

	c, err := ssh_config.Parse(fh)
	if err != nil {
		panic(err)
	}

	fh.Close()

	u, err := url.Parse(gitUri)
	if err != nil {
		panic(err)
	}

	host := c.FindByHostname(u.Host)
	idFile := host.GetParam("IdentityFile").Value()
	idFilePub := idFile + ".pub"

	ret, cred := git.NewCredSshKey("git", idFilePub, idFile, "")

	return git.ErrorCode(ret), &cred
}

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}

func normalizeGitUri(source string) (string, bool) {
	var gitUri string

	gitregexp := regexp.MustCompile("^(?:(https?|git|ssh)://|(git@))([^:|/]+)(?:/|:)([^/]+)/([^/\\.]+)(.git)$")
	u := gitregexp.FindStringSubmatch(source)

	if len(u) == 0 {
		return gitUri, false
	}

	var proto string
	if u[1] == "http" || u[1] == "https" || u[1] == "ssh" || u[1] == "git" {
		proto = u[1]
	} else {
		proto = "ssh"
	}

	gitUri = fmt.Sprintf("%s://%s%s/%s/%s%s", proto, u[2], u[3], u[4], u[5], u[6])

	return gitUri, true
}

func openGitRepo(source string) (*git.Repository, error) {
	var repo *git.Repository
	var err error

	gitUri, remote := normalizeGitUri(source)

	if remote {
		tmpdir, err := ioutil.TempDir("", "seekret")
		if err != nil {
			return nil, err
		}

		repo, err = git.Clone(gitUri, tmpdir, &git.CloneOptions{
			FetchOptions: &git.FetchOptions{
				RemoteCallbacks: git.RemoteCallbacks{
					CredentialsCallback:      credentialsCallback,
					CertificateCheckCallback: certificateCheckCallback,
				},
			},
		})
		if err != nil {
			return nil, err
		}
	} else {
		for {
			source, _ = filepath.Abs(source)
			repo, err = git.OpenRepository(source)
			if err == nil {
				break
			}
			if source == "/" {
				return nil,err
			}
			source = source + "/.."
		}
	}

	return repo, nil
}
