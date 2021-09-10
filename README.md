# README

This repository helps in identifying the permissions needed for the GitHub Actions.

It is not recommended to have write permissions to all, and it is recommended to follow this https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#permissions

But it is hard to determine the necessary permissions to existing actions, and the only way to do this is manually looking through the GitHub actions logs. https://docs.github.com/en/actions/reference/authentication-in-a-workflow#modifying-the-permissions-for-the-github_token

This tool automates getting the necessary permissions from the logs that can be set in the GitHub Action.

This is one of the checks for in `scorecard` to set correct permissions to avoid escalations.  https://github.com/ossf/scorecard/blob/main/docs/checks.md#token-permissions
## How to run

`export GITHUB_TOKEN=personalaccesstoken` https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token

`docker run -e GITHUB_AUTH_TOKEN=$GITHUB_TOKEN ghcr.io/naveensrinivasan/ghactionspermissions:latest ossf scorecard`

```
CodeQL .github/workflows/codeql-analysis.yml
permissions:
   contents: read
   metadata: read
   securityevents: write
   statuses: write
-------------------------------------------------------------------
build .github/workflows/main.yml
permissions:
   actions: write
   checks: write
   contents: write
   deployments: write
   discussions: write
   issues: write
   metadata: read
   packages: write
   pullrequests: write
   repositoryprojects: write
   securityevents: write
   statuses: write
-------------------------------------------------------------------
goreleaser .github/workflows/goreleaser.yaml
permissions:
   actions: write
   checks: write
   contents: write
   deployments: write
   discussions: write
   issues: write
   metadata: read
   packages: write
   pullrequests: write
   repositoryprojects: write
   securityevents: write
   statuses: write
-------------------------------------------------------------------
Close stale issues .github/workflows/stale.yml
permissions:
   actions: write
   checks: write
   contents: write
   deployments: write
   discussions: write
   issues: write
   metadata: read
   packages: write
   pullrequests: write
   repositoryprojects: write
   securityevents: write
   statuses: write
-------------------------------------------------------------------
Integration-tests .github/workflows/e2etests.yml
2021/09/10 21:09:38 Unable to fetch logs not a 200 status 502
Codescanning .github/workflows/codescan.yml
2021/09/10 21:09:39 Unable to fetch logs not a 200 status 410
Integration tests .github/workflows/integration.yml
permissions:
   actions: write
   checks: write
   contents: write
   deployments: write
   discussions: write
   issues: write
   metadata: read
   packages: write
   pullrequests: write
   repositoryprojects: write
   securityevents: write
   statuses: write
-------------------------------------------------------------------
Ok To Test .github/workflows/ok-to-test.yml
PR Verifier .github/workflows/verify.yml
permissions:
   actions: write
   checks: write
   contents: write
   deployments: write
   discussions: write
   issues: write
   metadata: read
   packages: write
   pullrequests: write
   repositoryprojects: write
   securityevents: write
   statuses: write
-------------------------------------------------------------------
```
