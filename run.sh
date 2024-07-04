docker image build -f Dockerfile -t ascii .
docker container run -p:8081:8081 --detach --name asciicontainer ascii
docker ps -a
echo To stop container, run "Docker stop <container_name>".

