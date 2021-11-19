#!/bin/bash
set -e

#Color to the people
RED='\x1B[0;31m'
CYAN='\x1B[0;36m'
GREEN='\x1B[0;32m'
NC='\x1B[0m'

#Variables
declare -a NODES=("0" "1" "2" "metachain")
GIT_HOME=~/observing-squad
STACK_FOLDER_MAINNET=~/MyObservingSquadMainnet
KEYS_FOLDER_MAINNET=${STACK_FOLDER_MAINNET}/keys
STACK_FOLDER_DEVNET=~/MyObservingSquadDevnet
KEYS_FOLDER_DEVNET=${STACK_FOLDER_DEVNET}/keys

OBSERVING_SQUAD_TAG=master

if [[ -f ~/.proxyDockerTag ]]; then
  PROXY_TAG=$(cat ~/.proxyDockerTag)
  echo ${PROXY_TAG}
fi

case "$1" in
'offline-devnet')
OFFLINE_ROSETTA_SERVER_PORT=$2
docker run -p "${OFFLINE_ROSETTA_SERVER_PORT}":8079 elrondnetwork/"${PROXY_TAG}" -offline -offline-config=./config/offline_config_devnet.toml
;;

'offline-mainnet')
OFFLINE_ROSETTA_SERVER_PORT=$2
docker run -p "${OFFLINE_ROSETTA_SERVER_PORT}":8079 elrondnetwork/"${PROXY_TAG}" -offline -offline-config=./config/offline_config_mainnet.toml
;;
'mainnet')

#Create the folder structure for the observer stack
mkdir -p ${STACK_FOLDER_MAINNET}/{proxy,node-0,node-1,node-2,node-metachain}/{config,logs}
mkdir -p ${STACK_FOLDER_MAINNET}/{node-0,node-1,node-2,node-metachain}/db
mkdir -p ${KEYS_FOLDER_MAINNET}

#Clone the repo and cd there
if [ -d "$GIT_HOME" ]; then sudo rm -rf $GIT_HOME; fi
git clone -b ${OBSERVING_SQUAD_TAG} https://github.com/ElrondNetwork/observing-squad.git $GIT_HOME
cd $GIT_HOME/rosetta-mainnet

if [[ -f ~/.proxyDockerTag ]]; then
   echo \ >> .env
  sed -i '/PROXY_TAG/d' .env
  echo PROXY_TAG=${PROXY_TAG} >> .env
fi

if [[ "$2" ]]; then
  sed -i '/PROXY_PORT/d' .env
  echo PROXY_PORT="$2" >> .env
fi;

#Generate Keys and place them in their respective folders

for OBSERVER_MAINNET in "${NODES[@]}"
do
   echo -e "${GREEN}--> Generating key for observer ${CYAN}$OBSERVER_MAINNET${GREEN}...${NC}"
   docker run --rm --mount type=bind,source=${KEYS_FOLDER_MAINNET},destination=/keys --workdir /keys elrondnetwork/elrond-go-keygenerator:latest && sudo chown $(whoami) ${KEYS_FOLDER_MAINNET}/validatorKey.pem && mv ${KEYS_FOLDER_MAINNET}/validatorKey.pem ${STACK_FOLDER_MAINNET}/node-$OBSERVER_MAINNET/config/observerKey_$OBSERVER_MAINNET.pem
done

#Say what has been started
echo "mainnet" > ~/.squadlocation

#Start the stack
echo -e
echo -e "${GREEN}--> Starting the Observer+Proxy MainNet Stack...${NC}"
echo -e
docker-compose --env-file .env up -d
;;

'devnet')
#Create the folder structure for the observer stack
mkdir -p ${STACK_FOLDER_DEVNET}/{proxy,node-0,node-1,node-2,node-metachain}/{config,logs}
mkdir -p ${STACK_FOLDER_DEVNET}/{node-0,node-1,node-2,node-metachain}/db
mkdir -p ${KEYS_FOLDER_DEVNET}

#Clone the repo and cd there
if [ -d "$GIT_HOME" ]; then sudo rm -rf $GIT_HOME; fi
git clone -b ${OBSERVING_SQUAD_TAG} https://github.com/ElrondNetwork/observing-squad.git $GIT_HOME
cd $GIT_HOME/rosetta-devnet

if [[ -f ~/.proxyDockerTag ]]; then
  echo \ >> .env
  sed -i '/PROXY_TAG/d' .env
  echo PROXY_TAG="${PROXY_TAG}" >> .env
fi

if [[ "$2" ]]; then
  sed -i '/PROXY_PORT/d' .env
  echo PROXY_PORT="$2" >> .env
fi;

#Generate Keys and place them in their respective folders

for OBSERVER_DEVNET in "${NODES[@]}"
do
   echo -e "${GREEN}--> Generating key for observer ${CYAN}$OBSERVER_DEVNET${GREEN}...${NC}"
   docker run --rm --mount type=bind,source=${KEYS_FOLDER_DEVNET},destination=/keys --workdir /keys elrondnetwork/elrond-go-keygenerator:latest && sudo chown $(whoami) ${KEYS_FOLDER_DEVNET}/validatorKey.pem && mv ${KEYS_FOLDER_DEVNET}/validatorKey.pem ${STACK_FOLDER_DEVNET}/node-$OBSERVER_DEVNET/config/observerKey_$OBSERVER_DEVNET.pem
done

#Say what has been started
echo "devnet" > ~/.squadlocation

#Start the stack
echo -e
echo -e "${GREEN}--> Starting the Observer+Proxy DEVNET Stack...${NC}"
echo -e
docker-compose --env-file .env up -d
;;

'stop')
read -p "Do you want to stop the stack ? [y/n] :" STOP

if [[ $STOP = y ]] ; then
                        if grep -Fxq "mainnet" ~/.squadlocation
                              then
                                 echo -e
                                 echo -e "${RED}--> Stopping the Observer+Proxy Stack...${NC}"
                                 echo -e
                                 cd $GIT_HOME/rosetta-mainnet && docker-compose down
                              else
                                 echo -e
                                 echo -e "${RED}--> Stopping the Observer+Proxy Stack...${NC}"
                                 echo -e
                                 cd $GIT_HOME/rosetta-devnet && docker-compose down
                        fi
            else
              echo -e
              echo -e "${GREEN}--> Ok...ignoring the command...${NC}"
              echo -e
fi

;;

*)
  echo "Usage: Missing parameter ! [mainnet|devnet|stop]"
;;

esac
