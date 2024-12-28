# Build an image named vue_helper using the Setup.Dockerfile
# The build args manage permissions when executing commands from inside the container
#
docker build -f ./dockerfiles/Dev.Dockerfile -t vue_app:dev vue_app
