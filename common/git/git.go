package git

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-ini/ini"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/common"
)

var (
	codeCommitHTTPRegex = regexp.MustCompile(`^https?://git-codecommit\.(.+)\.amazonaws.com/v1/repos/(.+)$`)
	codeCommitSSHRegex  = regexp.MustCompile(`ssh://git-codecommit\.(.+)\.amazonaws.com/v1/repos/(.+)$`)
	githubHTTPRegex     = regexp.MustCompile(`^https?://.*github.com.*/(.+)/(.+?)(?:.git)?$`)
	githubSSHRegex      = regexp.MustCompile(`github.com[:/](.+)/(.+?)(?:.git)?$`)

	cloneLock sync.Mutex

	ErrShortRef = errors.New("short SHA references are not supported")
	ErrNoRepo   = errors.New("unable to find git repo")
)

type Error struct {
	err    error
	commit string
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Commit() string {
	return e.commit
}

// FindGitRevision get the current git revision
func FindGitRevision(ctx context.Context, file string) (shortSha string, sha string, err error) {
	logger := log.Ctx(ctx)
	gitDir, err := findGitDirectory(file)
	if err != nil {
		return "", "", err
	}

	bts, err := os.ReadFile(filepath.Join(gitDir, "HEAD"))
	if err != nil {
		return "", "", err
	}

	var ref = strings.TrimSpace(strings.TrimPrefix(string(bts), "ref:"))
	var refBuf []byte
	if strings.HasPrefix(ref, "refs/") {
		// load commitid ref
		refBuf, err = os.ReadFile(filepath.Join(gitDir, ref))
		if err != nil {
			return "", "", err
		}
	} else {
		refBuf = []byte(ref)
	}

	logger.Info().Msgf("Found revision: %s", refBuf)
	return string(refBuf[:7]), strings.TrimSpace(string(refBuf)), nil
}

// FindGitRef get the current git ref
func FindGitRef(ctx context.Context, file string) (string, error) {
	logger := log.Ctx(ctx)
	gitDir, err := findGitDirectory(file)
	if err != nil {
		return "", err
	}
	logger.Debug().Msgf("Loading revision from git directory '%s'", gitDir)

	_, ref, err := FindGitRevision(ctx, file)
	if err != nil {
		return "", err
	}

	logger.Debug().Msgf("HEAD points to '%s'", ref)

	// Prefer the git library to iterate over the references and find a matching tag or branch.
	var refTag = ""
	var refBranch = ""
	r, err := git.PlainOpen(filepath.Join(gitDir, ".."))
	if err == nil {
		iter, err := r.References()
		if err == nil {
			for {
				r, err := iter.Next()
				if r == nil || err != nil {
					break
				}
				// logger.Debugf("Reference: name=%s sha=%s", r.Name().String(), r.Hash().String())
				if r.Hash().String() == ref {
					if r.Name().IsTag() {
						refTag = r.Name().String()
					}
					if r.Name().IsBranch() {
						refBranch = r.Name().String()
					}
				}
			}
			iter.Close()
		}
	}
	if refTag != "" {
		return refTag, nil
	}
	if refBranch != "" {
		return refBranch, nil
	}

	// If the above doesn't work, fall back to the old way

	// try tags first
	tag, err := findGitPrettyRef(ctx, ref, gitDir, "refs/tags")
	if err != nil || tag != "" {
		return tag, err
	}
	// and then branches
	return findGitPrettyRef(ctx, ref, gitDir, "refs/heads")
}

func findGitPrettyRef(ctx context.Context, head, root, sub string) (string, error) {
	logger := log.Ctx(ctx)
	var name string
	var err = filepath.Walk(filepath.Join(root, sub), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if name != "" || info.IsDir() {
			return nil
		}
		var bts []byte
		if bts, err = os.ReadFile(path); err != nil {
			return err
		}
		var pointsTo = strings.TrimSpace(string(bts))
		if head == pointsTo {
			// On Windows paths are separated with backslash character so they should be replaced to provide proper git refs format
			name = strings.TrimPrefix(strings.ReplaceAll(strings.Replace(path, root, "", 1), `\`, `/`), "/")
			logger.Debug().Msgf("HEAD matches %s", name)
		}
		return nil
	})
	return name, err
}

// FindGithubRepo get the repo
func FindGithubRepo(ctx context.Context, file, githubInstance, remoteName string) (string, error) {
	if remoteName == "" {
		remoteName = "origin"
	}

	url, err := findGitRemoteURL(ctx, file, remoteName)
	if err != nil {
		return "", err
	}
	_, slug, err := findGitSlug(url, githubInstance)
	return slug, err
}

func findGitRemoteURL(ctx context.Context, file, remoteName string) (string, error) {
	gitDir, err := findGitDirectory(file)
	if err != nil {
		return "", err
	}
	common.Logger(ctx).Debugf("Loading slug from git directory '%s'", gitDir)

	gitconfig, err := ini.InsensitiveLoad(fmt.Sprintf("%s/config", gitDir))
	if err != nil {
		return "", err
	}
	remote, err := gitconfig.GetSection(fmt.Sprintf(`remote "%s"`, remoteName))
	if err != nil {
		return "", err
	}
	urlKey, err := remote.GetKey("url")
	if err != nil {
		return "", err
	}
	url := urlKey.String()
	return url, nil
}

func findGitSlug(url string, githubInstance string) (string, string, error) {
	if matches := codeCommitHTTPRegex.FindStringSubmatch(url); matches != nil {
		return "CodeCommit", matches[2], nil
	} else if matches := codeCommitSSHRegex.FindStringSubmatch(url); matches != nil {
		return "CodeCommit", matches[2], nil
	} else if matches := githubHTTPRegex.FindStringSubmatch(url); matches != nil {
		return "GitHub", fmt.Sprintf("%s/%s", matches[1], matches[2]), nil
	} else if matches := githubSSHRegex.FindStringSubmatch(url); matches != nil {
		return "GitHub", fmt.Sprintf("%s/%s", matches[1], matches[2]), nil
	} else if githubInstance != "github.com" {
		gheHTTPRegex := regexp.MustCompile(fmt.Sprintf(`^https?://%s/(.+)/(.+?)(?:.git)?$`, githubInstance))
		gheSSHRegex := regexp.MustCompile(fmt.Sprintf(`%s[:/](.+)/(.+?)(?:.git)?$`, githubInstance))
		if matches := gheHTTPRegex.FindStringSubmatch(url); matches != nil {
			return "GitHubEnterprise", fmt.Sprintf("%s/%s", matches[1], matches[2]), nil
		} else if matches := gheSSHRegex.FindStringSubmatch(url); matches != nil {
			return "GitHubEnterprise", fmt.Sprintf("%s/%s", matches[1], matches[2]), nil
		}
	}
	return "", url, nil
}

func findGitDirectory(fromFile string) (string, error) {
	absPath, err := filepath.Abs(fromFile)
	if err != nil {
		return "", err
	}

	fi, err := os.Stat(absPath)
	if err != nil {
		return "", err
	}

	var dir string
	if fi.Mode().IsDir() {
		dir = absPath
	} else {
		dir = filepath.Dir(absPath)
	}

	gitPath := filepath.Join(dir, ".git")
	fi, err = os.Stat(gitPath)
	if err == nil && fi.Mode().IsDir() {
		return gitPath, nil
	} else if dir == "/" || dir == "C:\\" || dir == "c:\\" {
		return "", &Error{err: ErrNoRepo}
	}

	return findGitDirectory(filepath.Dir(dir))
}

type GitCloneConfig struct {
	URL   string
	Ref   string
	Dir   string
	Token string
}

// CloneIfRequired ...
func CloneIfRequired(ctx context.Context, refName plumbing.ReferenceName, input GitCloneConfig, lg *common.LoggerWrapper) (*git.Repository, error) {
	r, err := git.PlainOpen(input.Dir)
	if err != nil {
		cloneOptions := git.CloneOptions{
			URL:      input.URL,
			Progress: lg,
		}
		if input.Token != "" {
			cloneOptions.Auth = &http.BasicAuth{
				Username: "token",
				Password: input.Token,
			}
		}

		r, err = git.PlainCloneContext(ctx, input.Dir, false, &cloneOptions)
		if err != nil {
			lg.Errorf("Unable to clone %v %s: %v", input.URL, refName, err)
			return nil, err
		}

		if err = os.Chmod(input.Dir, 0755); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func gitOptions(token string) (fetchOptions git.FetchOptions, pullOptions git.PullOptions) {
	fetchOptions.RefSpecs = []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"}
	pullOptions.Force = true

	if token != "" {
		auth := &http.BasicAuth{
			Username: "token",
			Password: token,
		}
		fetchOptions.Auth = auth
		pullOptions.Auth = auth
	}

	return fetchOptions, pullOptions
}

// NewGitCloneExecutor creates an executor to clone git repos
//
//nolint:gocyclo
func Clone(ctx context.Context, input GitCloneConfig) error {

	logger := common.Logger(ctx)
	logger.Infof("  \u2601  git clone '%s' # ref=%s", input.URL, input.Ref)
	logger.Debugf("  cloning %s to %s", input.URL, input.Dir)

	cloneLock.Lock()
	defer cloneLock.Unlock()

	refName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", input.Ref))
	r, err := CloneIfRequired(ctx, refName, input, logger)
	if err != nil {
		return err
	}

	// fetch latest changes
	fetchOptions, pullOptions := gitOptions(input.Token)

	err = r.Fetch(&fetchOptions)
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	var hash *plumbing.Hash
	rev := plumbing.Revision(input.Ref)
	if hash, err = r.ResolveRevision(rev); err != nil {
		logger.Errorf("Unable to resolve %s: %v", input.Ref, err)
	}

	if hash.String() != input.Ref && strings.HasPrefix(hash.String(), input.Ref) {
		return &Error{
			err:    ErrShortRef,
			commit: hash.String(),
		}
	}

	// At this point we need to know if it's a tag or a branch
	// And the easiest way to do it is duck typing
	//
	// If err is nil, it's a tag so let's proceed with that hash like we would if
	// it was a sha
	refType := "tag"
	rev = plumbing.Revision(path.Join("refs", "tags", input.Ref))
	if _, err := r.Tag(input.Ref); errors.Is(err, git.ErrTagNotFound) {
		rName := plumbing.ReferenceName(path.Join("refs", "remotes", "origin", input.Ref))
		if _, err := r.Reference(rName, false); errors.Is(err, plumbing.ErrReferenceNotFound) {
			refType = "sha"
			rev = plumbing.Revision(input.Ref)
		} else {
			refType = "branch"
			rev = plumbing.Revision(rName)
		}
	}

	if hash, err = r.ResolveRevision(rev); err != nil {
		logger.Errorf("Unable to resolve %s: %v", input.Ref, err)
		return err
	}

	var w *git.Worktree
	if w, err = r.Worktree(); err != nil {
		return err
	}

	// If the hash resolved doesn't match the ref provided in a workflow then we're
	// using a branch or tag ref, not a sha
	//
	// Repos on disk point to commit hashes, and need to checkout input.Ref before
	// we try and pull down any changes
	if hash.String() != input.Ref && refType == "branch" {
		logger.Debugf("Provided ref is not a sha. Checking out branch before pulling changes")
		sourceRef := plumbing.ReferenceName(path.Join("refs", "remotes", "origin", input.Ref))
		if err = w.Checkout(&git.CheckoutOptions{
			Branch: sourceRef,
			Force:  true,
		}); err != nil {
			logger.Errorf("Unable to checkout %s: %v", sourceRef, err)
			return err
		}
	}

	if err = w.Pull(&pullOptions); err != nil && err != git.NoErrAlreadyUpToDate {
		logger.Debugf("Unable to pull %s: %v", refName, err)
	}
	logger.Debugf("Cloned %s to %s", input.URL, input.Dir)

	if hash.String() != input.Ref && refType == "branch" {
		logger.Debugf("Provided ref is not a sha. Updating branch ref after pull")
		if hash, err = r.ResolveRevision(rev); err != nil {
			logger.Errorf("Unable to resolve %s: %v", input.Ref, err)
			return err
		}
	}
	if err = w.Checkout(&git.CheckoutOptions{
		Hash:  *hash,
		Force: true,
	}); err != nil {
		logger.Errorf("Unable to checkout %s: %v", *hash, err)
		return err
	}

	if err = w.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: *hash,
	}); err != nil {
		logger.Errorf("Unable to reset to %s: %v", hash.String(), err)
		return err
	}

	logger.Debugf("Checked out %s", input.Ref)
	return nil
}
