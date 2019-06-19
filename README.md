# DevOpstic

![dashboard](doc/devopstic_demo.png "DevOpstic dashboard")

## Overview
DevOpstic allows you to store and display information related to your DevOps pipelines.

It's composed of:
- a REST api to receive events related to commit transition and incidents 
- a Postgres database for storage
- a Grafana instance to display a dashboards

## Metrics

The goal is to be able to track the following metrics, as identified in "Accelerate, Building and Scaling High Performing Technology Organization" by N. Forsgren, J. Humble & G. Kim:

- **Deployment Frequency**
- **Lead Time**
- **Mean Time To Restore (MTTR)**
- **Change Failure Rate**

**Deployment Frequency** 

Number of deployment per period of time

**Lead Time**

Time between a commit is submitted for review and the time it gets deployed.

It can be decomposed in two:
- **Review Lead Time**: from commit submitted for review to approved
- **Deployment Lead Time**: from approved to deployed successfully

**Mean Time To Restore**

Time between when an incident (failure of service) is detected to the time it is resolved.

**Change Failure Rate**

Ratio of failed deployment over total number of deployments.

*References:*
- https://stelligent.com/2018/12/21/measuring-devops-success-with-four-key-metrics/
- https://cloudplatformonline.com/2018-state-of-devops.html
- https://www.amazon.com/Accelerate-Software-Performing-Technology-Organizations-ebook/dp/B07B9F83WM

## Installation

- Kubernetes with Helm

`helm install https://github.com/splisson/devopstic/helm`

## Documentation
- [Domain](doc/concepts.md)

## License
- [MIT License](license.md)