package function

import (
	"crypto/subtle"
	"io/ioutil"
	"net/http"

	"code.gitea.io/sdk/gitea"
	"github.com/tidwall/gjson"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("HTTP Method Must be POST"))
		return
	}

	urlSecret := r.Header.Get("X-Buildkite-Token")
	envSecret, _ := getAPISecret("buildkite-secret")
	if !secureCompare(urlSecret, string(envSecret)) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Secret Validation failed"))
		return
	}

	orgSlug := queryString.Get("org_slug")
	if len(orgSlug) < 1 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Org Slug not defined"))
		return
	}
	repo := queryString.Get("repo")
	if len(repo) < 1 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Pipeline not defined"))
		return
	}

	var input []byte

	if r.Body == nil {
		// TODO: no json passed
		return
	}

	defer r.Body.Close()

	input, _ = ioutil.ReadAll(r.Body)

	status := gitea.CreateStatusOption{
		State:       getStateClean(gjson.Get(string(input), "build.state").String()),
		TargetURL:   gjson.Get(string(input), "build.web_url").String(),
		Description: getDescription(gjson.Get(string(input), "build.state").String()),
		Context:     "ci/buildkite",
	}

	apiBase, _ := getAPISecret("gitea-api-base")
	giteaToken, _ := getAPISecret("gitea-token")
	giteaClient, _ := gitea.NewClient(string(apiBase), gitea.SetToken(string(giteaToken)))

	_, _, err := giteaClient.CreateStatus(orgSlug, repo, gjson.Get(string(input), "build.commit").String(), status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Build Status was sent to Gitea"))
}

func getStateClean(state string) gitea.StatusState {
	switch state {
	case "pending":
		fallthrough
	case "running":
		fallthrough
	case "scheduled":
		return gitea.StatusPending
	case "failure":
		fallthrough
	case "failed":
		fallthrough
	case "declined":
		fallthrough
	case "canceled":
		return gitea.StatusFailure
	case "blocked":
		fallthrough
	case "success":
		fallthrough
	case "passed":
		return gitea.StatusSuccess
	case "warning":
		fallthrough
	default:
		return gitea.StatusWarning
	}
}

func getDescription(state string) string {
	switch state {
	case "pending":
		return "the build is pending"
	case "running":
		return "the build is running"
	case "the build is scheduled":
		return "pending"
	case "failure":
		return "the build failed"
	case "failed":
		return "the build failed"
	case "declined":
		return "the build was rejected"
	case "canceled":
		return "the build canceled"
	case "blocked":
		return "the build is pending approval"
	case "success":
		return "the build was successful"
	case "passed":
		return "the build has passed"
	case "warning":
		return "there has been a warning with the build"
	default:
		return "unknown"
	}
}

// function taken from https://docs.openfaas.com/reference/secrets/#use-the-secret-in-your-function
func getAPISecret(secretName string) (secretBytes []byte, err error) {
	// read from the openfaas secrets folder
	secretBytes, err = ioutil.ReadFile("/var/openfaas/secrets/" + secretName)
	if err != nil {
		// read from the original location for backwards compatibility with openfaas <= 0.8.2
		secretBytes, err = ioutil.ReadFile("/run/secrets/" + secretName)
	}

	return secretBytes, err
}

// function taken from https://play.golang.org/p/NU5uTaB-sp
func secureCompare(given string, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	} else {
		/* Securely compare actual to itself to keep constant time, but always return false */
		return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
	}
}
