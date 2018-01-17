# convox-off-cluster-builder
Tool to create Convox Builds off the Convox cluster. 


# usage

This is how we use it to do our build in CI and then import the build to the rack. (Using CirceCI in this case)

```makefile
// Makefile
./convox-build-off-cluster -app=service_name_here -compose-file=./docker-compose.convox.yml -description=${CIRCLE_BRANCH} -build-id=${CIRCLE_SHA1}
convox builds import -f service_name_here-${CIRCLE_SHA1}.tgz --rack rack_name_here

```
