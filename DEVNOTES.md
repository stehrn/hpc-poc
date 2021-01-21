# Ad hoc notes


https://console.cloud.google.com

## Articles 
* https://github.com/GoogleCloudPlatform/golang-samples
* [kubernes article this was inspired by](https://kubernetes.io/docs/tasks/job/fine-parallel-processing-work-queue/)
* https://blog.meain.io/2019/accessing-kubernetes-api-from-pod/
* https://cloud.google.com/pubsub/docs/quickstart-cli
* https://cloud.google.com/appengine/docs/flexible/go/writing-and-responding-to-pub-sub-messages
* https://github.com/googleapis/google-cloud-go/blob/master/pubsub/example_test.go
* https://pkg.go.dev/cloud.google.com/go/pubsub#example-Client.CreateSubscription
* https://cloud.google.com/run/docs/tutorials/pubsub
* https://cloud.google.com/kubernetes-engine/docs/tutorials/authenticating-to-cloud-platform
* Argo looks interesting : https://github.com/argoproj/argo/blob/master/workflow/controller/workflowpod.go

## Go specific build
https://cloud.google.com/cloud-build/docs/building/build-go

Pushes go binary to storage bucket, not used for container



gcloud pubsub subscriptions create tmp-nik-sub --topic projects/hpc-poc/topics/hpc-poc-bu3-topic  
projects/hpc-poc/subscriptions/tmp-nik-sub



func analyzeSentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeSentimentResponse, error) {
        return client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
                Document: &languagepb.Document{
                        Source: &languagepb.Document_Content{
                                Content: text,
                        },
                        Type: languagepb.Document_PLAIN_TEXT,
                },
        })
}

gcloud ml language analyze-entity-sentiment --content="I love R&B music. Marvin Gaye is the best. 'What's Going On' is one of my favorite songs. It was so sad when Marvin Gaye died."

https://github.com/mchmarny/tfeel