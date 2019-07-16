# How we use DevOpstic

Here is an example of how we use Devopstic.

## Pull requests

1. When a Pull Request is created, the Github Webhook is configured to  POST a request to
/github/events endpoint
2. When a Pull Request is approved and merged, same thing the Github Webhook invokes the github/events endpoint

## Deployments

Whenever we deploy a component, our Jenkins CI/CD jobs POST a request to /events endpoint.

## Failure

When a failure is detected by our Datadog monitoring, a POST request to /incidents endpoint is sent.