#!/usr/bin/env python3

import requests
import sys

bk_token = sys.argv[1]
env_secret = sys.argv[2]
branch = sys.argv[3]
commit = sys.argv[4]
username = sys.argv[5]
email = sys.argv[6]
bk_org_slug = sys.argv[7]
bk_pipeline = sys.argv[8]
message = sys.argv[9]

clean_branch = branch.split("/")

r = requests.post(("https://api.buildkite.com/v2/organizations/%s/pipelines/%s/builds" % (bk_org_slug, bk_pipeline) ),
	json={"commit":commit, "branch":clean_branch[-1], "message": message, "author":{ "name":username,"email":email }},
	headers={'Content-type': 'application/json', 'Accept': 'text/plain', 'Authorization': "Bearer %s" % bk_token })
