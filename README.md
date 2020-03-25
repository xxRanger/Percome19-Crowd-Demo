## Introduction

This demo is for paper **Aggregating Crowd Wisdom via Blockchain: A Private, Correct, and Robust Realization**, a system for secure aggregation. 

**The system only implement partial ideas of the original paper. If you are interested, please read the full paper**

The system includes three parties, data consumer, service provider and data providers. A data consumer can use this system to get an sum aggregation securely from data providers with help of service provider. The role of service provider is to aggregate encrypted data from data providers and use bulletproof to exclude invalid data, then use homomorphic addition to get the sum and return it to data consumer. 

A smart contract is deployed at Kovan testnet.
You can see from [here](https://kovan.etherscan.io/address/0x2a152cad883162011dcdf696a4861edc81619e14)

## Each stage of the system.

### Solicitation

Data consumer submit a task on blockchain.

### Registration

Service provider and data providers register on blockchain.

### Submission

Data providers submit encrypted data (a number at this demo) to blockchain. Note that data providers could be honest or malicious. 

### Aggregation

Service provider aggregate encrypted data from the blockchain. Use bulletproof to exclude invalid data (valid data is in a given range). Then use homomorphic addition to get the sum. Return sum back to blockchain.

### Approvment

Data consumer collect sum and approve to continue. 

### Claim

Service provider and data providers claim their reward.

## Each part of this project

### back-end

Provide core function for dealing with front-end and blockchain. 

### front-end

Subpath /admin: An admin page providing friendly interface for data consumer and service provider to manage crowdsourcing task.
Rootpath /: Provide friendly interface for data providers to submit their data and get their reward.

### client

Robots to simulate behavior of data consumer. You can set the number of data consumers to exploring this demo.

## Usage(As a test at local)

```bash
mkdir -p go-project/src/github.com/xxRanger
cd go-project/src/github.com/xxRanger
git clone https://github.com/xxRanger/Percome19-Crowd-Demo.git
```

### Run blockchain testnet

Follow the instruction to install ethereum https://geth.ethereum.org/docs/install-and-build/installing-geth

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

The default password for admin is conggroup, then click submit to continue

## Demo run flow

### Stage solicitation
admin: solicit task, as an example, input number 100, and click solicit
### Stage registration
user: submit a number like an invalid number -5 and click submit, at same time 20 robots will also submit
### Stage submission
admin: If no user register any more, stop registration.
### Stage aggregation
admin: Click aggregate. Server will aggregate results from blockchain and put sum back.
### Stage approvment
admin: Click show statics, it will showes the statistics of submissions, you can see your invalid number -5. Then click back to continue
admin: Click approve to continue.
### Stage claim
1. admin: Click claim to claim reward.
2. User: Click claim to claim reward.
3. After all party claim, back to step stage solicitation.



