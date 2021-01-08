#!/bin/sh

set +eu

# Configure accordingly
PROJECT_NAME="CHANGEME"
CLOUD_FUNCTION_URL="CHANGEME"
SCHEDULE="0 22 * * 6"
SERVICE_ACCOUNT_NAME="gitlab-environments-gc"
SCHEDULE_NAME="gitlab-environments-gc"

gcloud iam service-accounts create gitlab-environments-gc

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
