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
STACK_FOLDER_MAINNET=~/MyObservingSquad
KEYS_FOLDER_MAINNET=${STACK_FOLDER_MAINNET}/keys
STACK_FOLDER_TESTNET=~/MyObservingSquadTestnet
KEYS_FOLDER_TESTNET=${STACK_FOLDER_TESTNET}/keys

if [[ -f ~/.proxyDockerTag ]]; then
  PROXY_TAG=$(cat ~/.proxyDockerTag)
fi

case "$1" in
'offline')
OFFLINE_ROSETTA_SERVER_PORT=$2
docker run -p "${OFFLINE_ROSETTA_SERVER_PORT}":8079 elrondnetwork/"${PROXY_TAG}" -offline
;;
'mainnet')

#Create the folder structure for the observer stack
mkdir -p ${STACK_FOLDER_MAINNET}/{proxy,node-0,node-1,node-2,node-metachain}/{config,logs}
mkdir -p ${STACK_FOLDER_MAINNET}/{node-0,node-1,node-2,node-metachain}/db
mkdir -p ${KEYS_FOLDER_MAINNET}

#Clone the repo and cd there
if [ -d "$GIT_HOME" ]; then sudo rm -rf $GIT_HOME; fi
git clone -b master https://github.com/ElrondNetwork/observing-squad.git $GIT_HOME
cd $GIT_HOME/rosetta-mainnet

if [[ -f ~/.proxyDockerTag ]]; then
  sed -i '/PROXY_TAG/d' .env
  echo PROXY_TAG=${PROXY_TAG} >> .env
fi

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

'testnet')
#Create the folder structure for the observer stack
mkdir -p ${STACK_FOLDER_TESTNET}/{proxy,node-0,node-1,node-2,node-metachain}/{config,logs}
mkdir -p ${STACK_FOLDER_TESTNET}/{node-0,node-1,node-2,node-metachain}/db
mkdir -p ${KEYS_FOLDER_TESTNET}

#Clone the repo and cd there
if [ -d "$GIT_HOME" ]; then sudo rm -rf $GIT_HOME; fi
git clone -b master https://github.com/ElrondNetwork/observing-squad.git $GIT_HOME
cd $GIT_HOME/rosetta-testnet

if [[ -f ~/.proxyDockerTag ]]; then
  sed -i '/PROXY_TAG/d' .env
  echo PROXY_TAG=${PROXY_TAG} >> .env
fi

#Generate Keys and place them in their respective folders

for OBSERVER_TESTNET in "${NODES[@]}"
do
   echo -e "${GREEN}--> Generating key for observer ${CYAN}$OBSERVER_TESTNET${GREEN}...${NC}"
   docker run --rm --mount type=bind,source=${KEYS_FOLDER_TESTNET},destination=/keys --workdir /keys elrondnetwork/elrond-go-keygenerator:latest && sudo chown $(whoami) ${KEYS_FOLDER_TESTNET}/validatorKey.pem && mv ${KEYS_FOLDER_TESTNET}/validatorKey.pem ${STACK_FOLDER_TESTNET}/node-$OBSERVER_TESTNET/config/observerKey_$OBSERVER_TESTNET.pem
done

#Say what has been started
echo "testnet" > ~/.squadlocation

#Start the stack
echo -e
echo -e "${GREEN}--> Starting the Observer+Proxy TestNet Stack...${NC}"
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
                                 cd $GIT_HOME/rosetta-testnet && docker-compose down
                        fi
            else
              echo -e
              echo -e "${GREEN}--> Ok...ignoring the command...${NC}"
              echo -e
fi

;;

*)
  echo "Usage: Missing parameter ! [mainnet|testnet|stop]"
;;

esac
