# Introduction

This demo is for paper **Aggregating Crowd Wisdom via Blockchain: A Private, Correct, and Robust Realization**. 

A data consumer can use this app to get an sum aggregation of data providers. The smart contract is deployed at Kovan testnet.
You can see from [here](https://kovan.etherscan.io/address/0x2a152cad883162011dcdf696a4861edc81619e14)

## Every part of this project

### back-end

Provide core function for dealing with front-end and blockchain. 

### front-end

Subpath /admin: An admin page providing friendly interface for data consumer and service provider  to manage crowdsourcing task.
Rootpath /: Provide friendly interface for data providers to submit their data and get their reward.

### client

A simulation of data consumer. You can set the number of data consumers to exploring this demo.

## Usage 

```bash
mkdir -p go-project/src/github.com/xxRanger
cd go-project/src/github.com/xxRanger
git clone https://github.com/xxRanger/Percome19-Crowd-Demo.git
```

### Run blockchain testnet

Follow the instruction to in stall ethereum https://geth.ethereum.org/docs/install-and-build/installing-geth

Under folder Percome19-Crowd-Demo/testnet
```bash
./start.sh 
```
websocket port: 0.0.0.0:8650   
rpc port 0.0.0.0:8540

### Deploy the smart contract 

install nodejs https://nodejs.org/en/ 

install truffle  

```bash
 npm i -g truffle@v5.0.5
```

Deploy smart contract

Under folder Percome19-Crowd-Demo/truffle

```bash
cd truffle
truffle deploy
```

### Start server

You need to install go https://golang.org/dl/ and godep https://github.com/golang/dep

set GOPATH to /path/to/go-project

Under folder Percome19-Crowd-Demo/back-end
```bash
export GOPATH=$(cd ../../../../../ ;pwd)
echo $GOPATH
```
Install dependencies, may take ten minitues(or more) 
```bash
dep ensure
```

run server
```bash
go run main.go -chainPort "ws://0.0.0.0:8650"
```

### Start fake client 

Install dependencies 

Under folder Percome19-Crowd-Demo/client

```bash
dep ensure
```

run fake client

```bash
go run main.go 
```

### Start front end

Assume you have installed nodejs when install truffle

install dependencies

Under folder Percome19-Crowd-Demo/front-end

```bash
npm i 
```

if error try 
```bash
sudo npm i
```

Then start front-end

```bash
npm run serve 
```

Then front-end is running at localhost:8080

The page for user is at root path / 

The page for admin is at /#admin   

The default password for admin is conggroup, then click submit

## Run flow

1. admin: solicit task, as an example, input number 100, and click solicit
2. user: submit a number, say invalid number -5 as an example and click submit, at same time 20 fake clients will also submit
3. admin: No user register any more, so stop registration.
4. admin: Click Aggregate. Server will aggregate results from blockchain and put sum back.
5. admin: Click Approve to continue.
6. admin: Click show statics, it will showes the statistics of submission, you can see your invalid number -5. Then click back to continue
7. admin: Click claim to claim reward.
8. User: click claim to claim reward.
9. After all party claim, back to step 1.



