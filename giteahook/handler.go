package function

import (
	"crypto/subtle"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/buildkite/go-buildkite/v2/buildkite"
	"github.com/tidwall/gjson"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("HTTP Method Must be POST"))
		return
	}

	urlSecret := queryString.Get("secret")
	envSecret, _ := getAPISecret("gitea-secret")
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
	pipeline := queryString.Get("pipeline")
	if len(pipeline) < 1 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Pipeline not defined"))
		return
	}

	buildkiteSecret, _ := getAPISecret("buildkite-token")
	buildkiteConfig, err := buildkite.NewTokenConfig(string(buildkiteSecret), false)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := buildkite.NewClient(buildkiteConfig.Client())

	var input []byte

	if r.Body == nil {
		// TODO: no json passed
		return
	}

	defer r.Body.Close()

	input, _ = ioutil.ReadAll(r.Body)

	build := buildkite.CreateBuild{
		Commit:  gjson.Get(string(input), "after").String(),
		Branch:  trimRef(gjson.Get(string(input), "ref").String()),
		Message: gjson.Get(string(input), "commits.0.message").String(),
		Author: buildkite.Author{
			Name:  gjson.Get(string(input), "pusher.login").String(),
			Email: gjson.Get(string(input), "pusher.email").String(),
		},
	}

	_, _, err = client.Builds.Create(orgSlug, pipeline, &build)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Build was sent to Bulidkite"))
}

// TrimRef returns ref without the path prefix.
func trimRef(ref string) string {
	ref = strings.TrimPrefix(ref, "refs/heads/")
	ref = strings.TrimPrefix(ref, "refs/tags/")
	return ref
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
