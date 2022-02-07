# github-actions-pubsub-emulator-example

GCP makes an emulator available for the Pub/Sub service. They do this by providing it as a gcloud component that you are meant to use 9https://cloud.google.com/pubsub/docs/emulator).

In GitHub Actions, when you have a use case where you want to run something in the background throughout your workflow, like running the Pub/Sub emulator for tests, you'd normally use services (https://docs.github.com/en/actions/using-containerized-services/about-service-containers). GitHub made this feature with this use case in mind.

Unfortunately, right now, the services feature of workflows does not support overriding the command used to start the image used for the service. It can only start images with their default entry points. Their docs includes a Postgres example where this works fine:

```yaml
    # Service containers to run with `runner-job`
    services:
      # Label used to access the service container
      postgres:
        # Docker Hub image
        image: postgres
        # Provide the password for postgres
        env:
          POSTGRES_PASSWORD: postgres
        # Set health checks to wait until postgres has started
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          # Maps tcp port 5432 on service container to the host
          - 5432:5432
```

This means that to use the Pub/Sub emulator in a workflow, we must manually start it in a step, keeping it running in the background throughout the rest of the workflow. At the end of the workflow, the runner will automatically shut it down.

Example:

```yaml
      - name: 'Start Pub/Sub emulator'
        run: |
          gcloud beta emulators pubsub start --host-port=0.0.0.0:8085 &
```

For the full example, see `.github/workflows/use_pubsub_emulator.yml` in this repo.
