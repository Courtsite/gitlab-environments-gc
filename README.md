# gitlab-environments-gc

🗑️ A simple Google Cloud Function in Go to clean up stale [Environments in GitLab](https://docs.gitlab.com/ee/ci/environments/) - it can be triggered regularly using [Cloud Scheduler](https://cloud.google.com/scheduler).

---

**What environments are considered "stale"?**

- Its name does not contain protected words like `main`, `master`, `production`, and `protect` (we are only targetting preview / review environments)
- The last deployment was updated over 2 weeks ago

Both constraints are currently hardcoded, but we are open to making them configurable - pull requests welcomed.

**How does this work?**

This function does not clean up the actual environment (e.g. Kubernetes resources), you must define `on_stop` to achieve this (https://docs.gitlab.com/ee/ci/yaml/README.html#environmenton_stop).

Behind the scenes, the function will STOP the environment, and DELETE it.

To avoid hitting rate limits, and also timing out, this function currently only handles up to 200 environments at a time. If you have lots of environments, consider running this function more regularly.

**Why is this needed?**

At the time of writing, you can use `auto_stop_in` to automatically "stop" environments.
But, there is currently no automated / quick way to clean up stopped environments or environments started without `auto_stop_in`.
This project is essentially a stop-gap to help you achieve that.


## Getting Started

### Prerequisites

- Ensure you have `gcloud` installed:
    - MacOS: `brew cask install google-cloud-sdk`
    - Others: https://cloud.google.com/sdk/gcloud
- Ensure you have authenticated with Google Cloud: `gcloud init`
- (Optional) Set your current working project: `gcloud config set project <project>`

### Deployment

1. Clone / download a copy of this repository
2. Copy `.env.sample.yaml` to `.env.yaml`, and modify the environment variables declared in the file
3. Run `./deploy.sh` _(recommendation: do not allow unauthenticated requests, see section on Cloud Scheduler below for more information)_

### Setting Up Cloud Scheduler

Optionally, if you would like to run the clean up regularly (cron), modify, and run `./create_scheduler.sh` to create the necessary resources.
Unlike `./deploy.sh` which you can run multiple times, you should only run `./create_scheduler.sh` once.
This script will also set-up a service account which will be used to securely invoke the Cloud Function, without exposing it to the internet (specifically, preventing unauthenticated requests).

Alternatively, you can also set-up Cloud Scheduler manually via the [web UI](https://console.cloud.google.com/cloudscheduler) or through _infrastructure-as-code_, e.g. [Terraform](https://registry.terraform.io/providers/hashicorp/google/latest/docs).
