# IAM
Information aobut service account and role create


# Application 
```
export PROJECT_NAME=hpc-poc
```

## Client 
```
gcloud iam roles create hpc.client --project=${PROJECT_NAME} --file=client_role.yaml
gcloud iam roles describe --project=${PROJECT_NAME} hpc.client
```

# Integration testing
Create service account:
```
gcloud iam service-accounts create hpc-integration-test --description="Integration test account" --display-name="Integration test"
export SERVICE_ACCOUNT=hpc-integration-test@${PROJECT_NAME}.iam.gserviceaccount.com
gcloud iam service-accounts describe ${SERVICE_ACCOUNT}
```
Create custom role:
```
gcloud iam roles create hpc.integration.test --project=${PROJECT_NAME} --file=integration_test.yaml
gcloud iam roles describe --project=hpc-poc hpc.client
```
Bind role to SA:
```
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=projects/${PROJECT_NAME}/roles/hpc.integration.test
```
Create key:
```
gcloud iam service-accounts keys create ${HOME}/integration_test_key.json --iam-account ${SERVICE_ACCOUNT}
```
Export key for use in integration tests:
```
export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/integration_test_key.json
```


# Reference
* https://cloud.google.com/iam/docs/creating-custom-roles