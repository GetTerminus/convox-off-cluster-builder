# convox-off-cluster-build
Tool to create Convox Builds off the Convox cluster.

## Enhancements:
The tool takes a docker-compose.convox.yml file and turns the `build:` section into an `image:` section to speed up subsequential builds.
To achieve this `convox-off-cluster-build` builds the project specified in the `build:` section, creates an image from it and uploads the image into an ECR's Repository.
It rewrites the docker-compose.convox.yml file and removes the `build:` section from it and substitutes it with an `image:` section that points to the ECR Repository and the correct image.

Subsequent builds will pick up the image and avoid a rebuild which is more time consuming than an image copy from AWS ECR to the destination.

The new `convox-build-off-cluster` tool is intended to be part of a chain of commands that builds and deploys services.
It will only become active if there is a `build:` section in the docker-compose.convox.yml file. If it sees an `image:` section for a given service, it simply skips it.

# Usage



`convox-off-cluster-build -h` for help

# Configuration Example (_.drone.yml_)
```
convox_ninja:
  image: getterminus/cci-build-golang:20180426
  environment:
	- AWS_REGION=us-east-1
  secrets:
	- source: build_convox_host
	  target: convox_host
	- source: build_convox_password
	  target: convox_password
	- source: build_aws_account
	  target: aws_account
	- source: build_aws_access_key_id
	  target: aws_access_key_id
	- source: build_aws_secret_access_key
	  target: aws_secret_access_key
	- source: build_repo
	  target: repo
  commands:
	- convox-build-off-cluster -app=<your app name> -description=${DRONE_COMMIT_BRANCH} -gitsha=${DRONE_COMMIT_SHA} -repo $${REPO}
	- mkdir scratch
	- mv ./docker-compose.convox.yml scratch/docker-compose.yml
	- cd ./scratch
	- convox build --app=<your app name>
```

## Set the AWS_REGION
```yaml
environment:
  - AWS_REGION=us-east-1
```

## Last, but not least, there is one more thing to take care of
### This step is a cane, but it is needed to get the RELEASEID.
This step prevents convox from actually building the service yet again.
```yaml
commands:
  - mkdir scratch
  - mv ./docker-compose.convox.yml scratch/docker-compose.yml
  - cd ./scratch
  - convox build --app=<your app name>
```
The `convox build...` command picks up the image that has been built and pushed by the previous `convox-build-off-cluster` tool and turns it into a release.
