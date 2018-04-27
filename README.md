# convox-off-build-off-cluster
Tool to create Convox Builds off the Convox cluster.


# Usage
This is how we use it to do our build in CI and then import the build to the rack. (Using CirceCI in this case)
The new convox-build-off-cluster tool has been overhauled and its functionality has been expanded.
For the new and improved version of the convox-build-off-cluster tool an adjusted version of the .drone.yml file is needed

# Configuration Example (_.drone.yml_)
```yaml
convox_ninja:
  image: getterminus/cci-build-golang:20180426  <-- latest image
  environment:
	- AWS_REGION=us-east-1
  secrets:
	- source: ninja_convox_host
	  target: convox_host
	- source: ninja_convox_password
	  target: convox_password
	- source: ninja_aws_account
	  target: aws_account
	- source: ninja_aws_access_key_id
	  target: aws_access_key_id
	- source: ninja_aws_secret_access_key
	  target: aws_secret_access_key
	- source: ninja_repo
	  target: repo
  commands:
	- convox-build-off-cluster -app=<your app name> -description=${DRONE_COMMIT_BRANCH} -gitsha=${DRONE_COMMIT_SHA} -repo $${REPO}
	- mkdir scratch
	- mv ./docker-compose.convox.yml scratch/docker-compose.yml
	- cd ./scratch
	- convox build --app=<your app name>
```

## There are several secrets that need to be configured in the drone "Secrets" from the drone UI
These settings will need to be provided by the SRE team
  * ninja_aws_account
  * ninja_aws_access_key_id
  * ninja_aws_secret_access_key
  * ninja_repo

## Additionally to these mandatory "secrets" settings you will need to add the region for your service by an environment variable
```yaml
environment:
  - AWS_REGION=us-east-1
```
Adjust for your region (as of now it's only 'us-east-1') as needed
## Last, but not least, there is one more thing to take care of
```yaml
commands:
  - mkdir scratch
  - mv ./docker-compose.convox.yml scratch/docker-compose.yml
  - cd ./scratch
  - convox build --app=<your app name>
```
### This step is a cane, but it is needed to get the RELEASEID.
This step prevents convox from actually building the service yet again.

It picks up the image that has been built and pushed by the previous **convox-build-off-cluster** tool. The **convox** tool picks up the image and turns it into a release.
