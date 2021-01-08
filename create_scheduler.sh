#!/bin/sh

set +eu

# WARNING: run this script only once when setting up for the first time!

# Configure accordingly
PROJECT_NAME="CHANGEME"  # This is your Google Cloud project name
CLOUD_FUNCTION_URL="CHANGEME"  # This is the HTTP(S) trigger URL returned after successfully running `deploy.sh`

# https://crontab.guru/#0_22_*_*_6
SCHEDULE="0 22 * * 6"

# You do not have to change these, unless you prefer different names.
SERVICE_ACCOUNT_NAME="gitlab-environments-gc"
SCHEDULE_NAME="gitlab-environments-gc"

gcloud iam service-accounts create "$SERVICE_ACCOUNT_NAME"

gcloud projects add-iam-policy-binding "$PROJECT_NAME" \
      --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_NAME.iam.gserviceaccount.com" \
      --role="roles/cloudfunctions.invoker"

gcloud scheduler jobs create http \
    "$SCHEDULE_NAME" \
    --uri="$CLOUD_FUNCTION_URL" \
    --schedule="$SCHEDULE" \
    --time-zone="Etc/UTC" \
    --http-method=GET \
    --max-retry-attempts=0 \
    --oidc-service-account-email="$SERVICE_ACCOUNT_NAME@$PROJECT_NAME.iam.gserviceaccount.com"
