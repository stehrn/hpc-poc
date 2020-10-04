

# Console
https://console.cloud.google.com

# Authentication
https://cloud.google.com/docs/authentication/getting-started
export GOOGLE_APPLICATION_CREDENTIALS="/Users/db/Downloads/GOOGLE-HPC-POC.json"

# Writing a function 
https://cloud.google.com/functions/docs/quickstart-console

#  Deploy
https://cloud.google.com/functions/docs/deploying/filesystem

gcloud init
github.com/stehrn/gcp/function.go

gcloud functions deploy local-deploy-function --entry-point function.go --runtime go113 --trigger-topic MY_TOPIC
