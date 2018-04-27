# convox-off-cluster-build
Tool to create Convox Builds off the Convox cluster.

## Enhancements:
The tool takes a `docker-compose.convox.yml` file and turns the `build:` section into an `image:` section to speed up subsequential builds.
To achieve this `convox-off-cluster-build` builds the project specified in the `build:` section, creates an image from it and uploads the image into an ECR's Repository.
It rewrites the `docker-compose.convox.yml` file and removes the `build:` section from it and substitutes it with an `image:` section that points to the ECR Repository and the correct image.

Note that a single ECR Repository can only host 1,000 images, so ensure that you have a Lifecycle Policy in place.
See the `Lifecycle Policy` tab on AWS ECR / Repositories / YourRespository for details.

Subsequent builds will pick up the image and avoid a rebuild which is more time consuming than an image copy from AWS ECR to the destination.

The new `convox-off-cluster-build` tool is intended to be part of a chain of commands that builds and deploys services.
It will only become active if there is a `build:` section in the `docker-compose.convox.yml` file. If it sees an `image:` section for a given service, it simply skips this service.

# Usage
## Getting Help
`convox-off-cluster-build -h` for help

## Easy invocation
`convox-off-cluster-build -app <your_application_name> -gitsha <add your current gitsha here> -region us-east-1 -defaultAccount 1234567890 -repo <your_repo>`

## This tricky little step prevents convox from actually building the service yet again
```sh
mkdir scratch
mv ./docker-compose.convox.yml scratch/docker-compose.yml
cd ./scratch
convox build --app=<your app name>
```
The `convox build...` command picks up the `image:` section that has been created by the previous `convox-off-cluster-build` tool and turns it into a release.
