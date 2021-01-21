# IAM
Details on service accounts and roles

# Application 
Set project name:
```
export PROJECT_NAME=hpc-poc
```

## Client 
Create service account:
```
gcloud iam service-accounts create hpc-client --description="HCP client account" --display-name="HCP client account"
export SERVICE_ACCOUNT=hpc-client@${PROJECT_NAME}.iam.gserviceaccount.com
gcloud iam service-accounts describe ${SERVICE_ACCOUNT}
```
Create custom role:
```
gcloud iam roles create hpc.client --project=${PROJECT_NAME} --file=client_role.yaml
gcloud iam roles describe --project=${PROJECT_NAME} hpc.client
```
Bind role to SA:
```
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=projects/${PROJECT_NAME}/roles/hpc.client
```
Create key:
```
gcloud iam service-accounts keys create ${HOME}/client_key.json --iam-account ${SERVICE_ACCOUNT}
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


xxx
## Summary

* Orchestrator:
  * subscribe [`pubsub.subscriber`]
  * view [`pubsub.viewer`] (to check subscription exists)
  * delete storage objects [`storage.objects.delete`]
  * publish [`roles/editor`]

* Engine:
  * read storage objects [`storage.objectViewer`]

* Monitor:
  * read storage objects [`storage.objectViewer`]


(see [authenticating-to-cloud-platform](https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform))

For now, roles are v broad and grouped together, so more work needed here
```
export SERVICE_ACCOUNT=gke-sub-acc@hpc-poc.iam.gserviceaccount.com

// create service sccount 
gcloud iam service-accounts create ${SERVICE_ACCOUNT} --description="GKE subscription account" --display-name="gke-subscription"
gcloud iam service-accounts list

// add roles
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/pubsub.subscriber 
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/pubsub.publisher
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/pubsub.viewer
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/storage.objectAdmin
gcloud projects add-iam-policy-binding ${PROJECT_NAME} --member=serviceAccount:${SERVICE_ACCOUNT} --role=roles/pubsub.editor


// list roles
gcloud projects get-iam-policy ${PROJECT_NAME} --flatten="bindings[].members" --format='table(bindings.role)' --filter="bindings.members:${SERVICE_ACCOUNT}"
```

* [storage iam-permissions](https://cloud.google.com/storage/docs/access-control/using-iam-permissions)

xxx