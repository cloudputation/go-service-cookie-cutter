# !/bin/bash

SERVICE_PORT="9595"
HOME_ITERATOR_DATA="/opt/iterator/data/"
DOCKER_ITERATOR_DATA="/var/lib/iterator"

rebuild() {
  for c in $(sudo docker ps -a | grep -i "exited" | awk '{print $1}')
  do
    sudo docker rm $c
  done

  for i in $(sudo docker images | grep "none" | awk '{print $3}')
  do
    sudo docker rmi $i
  done

  sudo chown -R 991:991 $HOME_ITERATOR_DATA

  sudo docker run -d \
    -p ${SERVICE_PORT}:${SERVICE_PORT} \
    -v ${HOME_ITERATOR_DATA}:${DOCKER_ITERATOR_DATA} \
    --name iterator iterator:0.0.2

  sudo docker logs --follow iterator
}


sudo docker stop iterator || true
sudo docker rm iterator || true
mkdir -p ./artifacts
cp ${HOME}/dev/artifacts/terraform/* ./artifacts

sudo rm ./build/iterator
sudo make docker-build && rebuild
