package function

import (
	"crypto/subtle"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/buildkite/go-buildkite/v2"
	"github.com/itchyny/gojq"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("HTTP Method Must be POST"))
		return
	}

	urlSecret := queryString.Get("secret")
	envSecret := getAPISecret("gitea-secret")
	if !secureCompare(urlSecret, envSecret) {
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

	var payload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buildkiteConfig, _ := buildkite.NewTokenConfig(getAPISecret("buildkite-secret"), false)

	client := buildkite.NewClient(buildkiteConfig.Client())

	build := buildkite.CreateBuild{
		Commit:  payload["after"],
		Branch:  payload["ref"],
		Message: payload["commits"][0]["message"],
		Author: buildkite.Author{
			Name:  payload["pusher"]["login"],
			Email: payload["pusher"]["email"],
		},
	}

	_, _, err = client.BuildService.Create(orgSlug, payload, build)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Build was sent to Bulidkite"))
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
