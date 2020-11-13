# Gitea/Buildkite Connector

OpenFaaS functions that connect Gitea and Bulidkite to trigger builds and report build statuses back

## Getting started

### Create Secrets

First you'll need to create the secrets that that once deployed the functions can start to work

* Buildkite API Key

Go to [Bulidkite Settings](https://buildkite.com/user/api-access-tokens) and create an API key that has `write_builds` REST scope

`faas secret create buildkite-token --from-literal "BUILDKITE_API_KEY"`

* Buildkite Webhook Secret

This is a shared secret that'll verify data sent to function. You can choose what to make this string, just make sure to make it secure

`faas secret create buildkite-secret --from-literal "BUILDKITE_WEBHOOK_SECRET"`

* Gitea API Key

Go to Gitea user application settings https://gitea.example.com/user/settings/applications and create an API token

`faas secret create gitea-token --from-literal "GITEA_API_KEY"`

* Gitea Webhook Secret

This is a shared secret that'll verify data sent to function. You can choose what to make this string, just make sure to make it secure

`faas secret create gitea-secret --from-literal "GITEA_WEBHOOK_SECRET"`

* Gitea API Base

This is the URL to the Gitea install so that API calls know where to call

`faas secret create gitea-api-base --from-literal "https://gitea.example.com"`


### Deploy your function

Just like any OpenFaaS function, just deploy the function using `faas deploy -f stack.yml`

### Setup Webhooks

Now that secrets have been made, you'll need to create webhooks in both Gitea and Buildkite

#### Create Buildkite Webhook

First go to the "Notification Services" page of your Org settings, create a "Webhook" with all of the "Build" events checked, and select the specific pipeline that is connected to your repo.

Add two querystring params to your function URL `org_slug=username_or_org_slug_from_gitea` and `repo=repo_slug_from_gitea`, and then use that as URL for webhook. It'll end up looking something like: `https://openfaas.example.com/function/buildkitehook?org_slug=repo_slug_from_gitea&repo=repo_slug_from_gitea`

Finally, add the secret you created above as the webhook token.

Now that the Buildkite webhook has been created, the status of your builds will be reported back to Gitea.

#### Create Buildkite Webhook

Go to the webhook settings section of repo settings and create a new "Gitea" webhook. Only select "Push Events" to trigger on.

Just like with the buildkite webhook, you'll need to add several querystring params to the function URL. `secret=GITEA_WEBHOOK_SECRET`, `org_slug=buildkite_org_slug`, and `pipeline=buildkite_pipeline_slug`, so it'll end up looking like `https://openfaas.example.com/function/giteahook?secret=GITEA_WEBHOOK_SECRET&org_slug=buildkite_org_slug&pipeline=buildkite_pipeline_slug`
