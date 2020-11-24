

A simple engine, it actually just reads value of `PAYLOAD` env variable, prints it out, and exits.

To build image:
```
go mod init
gcloud builds submit --tag gcr.io/hpc-poc/engine
```

Image is referenced in orchestrator deployment (`ENGINE_IMAGE`)

To test engine:
TODO