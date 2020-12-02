
Look into quotas: https://cloud.google.com/kubernetes-engine/quotas


github.com/googleapis/gnostic/OpenAPIv2




Flow:
* client 
  * write data to cloud storage bucket
  * publish message containing data location (bucket/object)
* orchestrator 
  * subscribes to topic
    * on message - create kubernetes Job, passing it location of cloud storage data
  * watches jobs
    * on job success - delete cloud storage object
* engine
   * reads cloud storage data, does something with it, exists

Web applications:
 * client - submit data 
 * monitor - view jobs and pods (extend to view pubsub details)