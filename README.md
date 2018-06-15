# Gitea/Buildkite Connector

This docker image hosts a webhook receiver for both Gitea and Buildkite webhooks to be able to use Buildkite CI to automatically buildu Gitea commits, and report status back to Gitea.

## Environment variables

* `GT_API_BASE`: If your Gitea API can be reached at `https://git.example.com/api` then this variable should be `https://git.example.com`. Note: No trailing slash
* `GT_TOKEN`: Gitea API Token that is used to authenticate against the Gitea API to report status back to Gitea. You can generate this from the applications page under user settings `https://git.example.com/user/settings/applications`
* `GT_SECRET`: Webhook secret used to ensure that webhook is coming from Gitea. You'll configure this when creating your webhook.
* `BK_TOKEN`: Buildkite API Token that is used to authenticate against the Buildkite API to create new builds. Setup a new token under api access tokens with write_build permissions at `https://buildkite.com/user/api-access-tokens`
* `BK_SECRET`: Webhook secret used to ensure that webhook is coming from Buildkite. Get this from when you setup a webhook with all the build events under organization > notification settings.

## Setup

Gather environment variables from above, and put fill them out in the below docker run command:

```
docker run \
-p 9000:9000 \
-e GT_SECRET=<FILL_WITH_REAL_VALUE> \
-e GT_TOKEN=<FILL_WITH_REAL_VALUE> \
-e GT_API_BASE=https://git.example.com \
-e BK_TOKEN=<FILL_WITH_REAL_VALUE> \
-e BK_SECRET=<FILL_WITH_REAL_VALUE> \
-d \
techknowlogick/gitea-buildkite:latest
```

When setting up your webhook on Gitea you'll need to pass org and pipeline in the querystring: example `https://git.example.com:9000/hooks/gitea-webhook?org_slug=organization_slug_example&pipeline=pipeline_example`

Your webhook from Buildkite will look like: `http://git.example.com:9000/hooks/buildkite-webhook`
