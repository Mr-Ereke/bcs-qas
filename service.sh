#!/usr/bin/bash

confMinor=$(pwd)/src/config/env.conf
. $confMinor

dockerBuild () {
  echo "Docker building..."
  echo "Docker container name " $CONTAINER_NAME
  docker build -f ./DockerfileLocal -t $CONTAINER_NAME .
  echo "Docker removing builder images..."
  docker rmi $(docker images --filter label=stage=builder -q)
}

if [ "$1" == 'start' ] || [ "$1" == 'restart' ]; then
  if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    if [ "$1" == 'restart' ]; then
      echo "Restarting..."
      docker container stop $CONTAINER_NAME
    else
      echo "Service is already running"
      exit
    fi
  fi

  if [ "$(docker ps -aq -f status=exited -f name=$CONTAINER_NAME)" ]; then
    # cleanup, if container exited/stopped
    echo "Removing container..."
    docker container rm $CONTAINER_NAME
  fi

  echo "Docker running..."
  echo $CONTAINER_NAME
  docker run -d --restart unless-stopped \
    --name=$CONTAINER_NAME \
    --network=host \
    --mount type=bind,src=/etc/hostname,dst=/etc/hostname \
    --log-opt max-size=30m --log-opt max-file=3 \
    -w /go $CONTAINER_NAME

  echo "Docker logging..."
  docker logs $CONTAINER_NAME -f >> $HOST_LOGS_DIRECTORY/$CONTAINER_NAME.log
elif [ "$1" == 'stop' ]; then
  echo "Docker stopping..."
  docker container stop $CONTAINER_NAME
elif [ "$1" == 'build' ]; then
  dockerBuild
elif [ "$1" == 'rebuild' ]; then
  echo "Docker rebuilding..."
  if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
      echo "Stopping $CONTAINER_NAME container..."
      docker container stop $CONTAINER_NAME

      echo "Removing $CONTAINER_NAME container..."
      docker container rm $CONTAINER_NAME
  fi

  if [ "$(docker images -q $CONTAINER_NAME)" ]; then
      echo "Removing $CONTAINER_NAME image..."
      docker rmi $(docker images -q $CONTAINER_NAME)
  fi

  echo "Waiting for delete images..."

  # Waits 2 seconds for delete images
  sleep 2s

  echo "Docker recreating..."

  dockerBuild
fi
