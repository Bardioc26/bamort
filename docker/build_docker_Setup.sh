# Build an image named vue_helper using the Setup.Dockerfile
# The build args manage permissions when executing commands from inside the container
#
docker build \
  --build-arg USER_ID=$(id -u) \
  --build-arg GROUP_ID=$(id -g) \
  -t vue_helper:slim - < ./dockerfiles/Setup.Dockerfile
