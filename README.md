# Incomplete Kaleido Project

##  Overview

This is currently incomplete but the high points are:

- Use go-redis and redis db as a persistent store.
- I found gofiber and am using that for the web server and rest API bits.  I started with go-kit early on and it worked, but it was more complex so I switched to gofiber.
- I'm doing all this on a local quorum enviroment stood up by the kaleido quorum-tools.  This is all running in a local Linux VM.
- Dev environment is vscode.
- I grabbed the ERC721 gameItem contract from Openzeppelin and am writing a rest API to award, trade and delete in game items.  This wasn't my original intent on a usecase, but it's where it is.
- I went ahead and precompiled the contract using solc and abigen although I was able to compile it by modifying the kaleido example to add the node path.

## Setup

1.  Get Docker and golang.  I'm using 1.17.4 of golang and whatever the latest docker is.
2.  Setup and deploy quorum-tools to get an ethereum environment.  The code currently assumes a node at localhost:22001.  Also make sure to get the chainid.  It's not 2018.
3.  Clone this repo
1.  run `nobuild.sh` to build a container with the go binary included (see below for details).
2.  run `docker-compose up` to stand up redis and the kec app `apptwo` and create the proper network.

Alternatively you could do the following:

1.  run `runit.sh` to deploy the container as apptwo.  (We don't talk about appone.  Too soon.)
1.  run `runredis.sh` to deploy redis.  Note that this also puts Redis and apptwo on a kec network so I can refer to redis by name.

That's it.  You can use curl to 8081 to verify the API is up.

### Details on the scripts

- buildit.sh: This shell script cleans up old containers and images then uses the `Dockerfile-buildit` dockerfile which builds the project in the container.  I don't use this because it's faster to build locally and just put the binary in the container.
- nobuild.sh: This shell script also cleans up old containers and images then uses `Dockerfile` which just copies the latest build binary into the container.  Much faster and smaller container, but of course you have to build locally first.
- runit.sh: Runs the app as apptwo.  Puts it on the proper docker network and exposes port 8081 on the host to the container.
- runredis.sh: Runs redis on the proper network.

### Details on the API

TODO: Swagger.

- `GET /` Simple check to ensure the app is up and running.  
ex: `curl http://localhost:8081`
- `POST /api/v1/user/id` Creates a user/transactor (represented by a private key) for Ethereum and stores the key in Redis.  
ex: `curl -X POST http://localhost:8081/api/v1/user/paul` will create a user paul.  
- `POST /api/v1/deployContract/id` Deploys the GameItem contract by a particular user specified by ID  
ex: `curl -X POST http://localhost:8081/api/v1/deployContract/paul` to deploy the GameItem contract as user paul.
- `GET /api/v1/testFunction` Does some random Ethereum stuff.  I use this to test new things before I add to an actual method.  
ex: `curl http://localhost:8081/testFunction`  I don't really remember what this returns but it does stuff.
- `POST /api/v1/awardItem/user/itemname/instance` Awards item itemname to specified user via contract at address instance.  
ex: `curl -X POST http://localhost:8081/api/v1/awardItem/paul/hammer/0x739d7B9E7552E972b1d1F73Acf938BcF34eCb9C1`
- `GET /api/v1/getItems` Gets a list of available items to be awarded.  Right now there's only a hammer because THOR!  
ex: `curl http://localhost:8081/api/v1/getItems`

Planned APIs

- `POST /api/v1/tradeItem/itemid` Trades an item from the owner to another user.
- `POST /api/v1/destroyItem/itemid` Destroys a users item.  Should only be able to be done by that user.
Also some Get APIs around viewing items, etc.
