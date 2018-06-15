#!/usr/bin/env python3

import requests
import sys

gt_token = sys.argv[1]
env_secret = sys.argv[2]
url_secret = sys.argv[3]
repo = sys.argv[4]
commit_id = sys.argv[5]
result_url = sys.argv[6]
event = sys.argv[7]
state = sys.argv[8]
api_base = sys.argv[9]

if env_secret != url_secret:
	print("secrets don't match")
	exit()

repo_1 = repo.split(":")[1]
repo_2 = repo_1.split("/")
owner = repo_2[0]
repo_clean = repo_2[1][:-4]

def get_state_clean(state):
    return {
        'pending': "pending",
        'running': "pending",
        'success': "success",
        'failure': "failure",
        'canceled': "failure",
        'blocked': "warning",
        'declined': "failure",
        'passed': "success",
        'scheduled': 'pending'
    }.get(state, 'warning')

def get_desc(state):
	return {
		'scheduled': 'the build is scheduled',
        'pending': "the build is pending",
        'running': "the build is running",
        'success': "the build was successful",
        'failure': "the build failed",
        'canceled': "the build canceled",
        'blocked': "the build is pending approval",
        'declined': "the build was rejected",
        'passed':  "the build has passed"
    }.get(state, 'unknown')

r = requests.post(("%s/api/v1/repos/%s/%s/statuses/%s" % (api_base, owner, repo_clean, commit_id)),
	json={"context":"ci/buildkite", "target_url":result_url, "state": get_state_clean(state), "description":get_desc(state)},
	headers={'Content-type': 'application/json', 'Accept': 'text/plain', 'Authorization': "token %s" % gt_token })
