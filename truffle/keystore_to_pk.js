const Wallet = require('ethereumjs-wallet'),
      fs = require('fs');


const utcFile = "/Users/anxin/Desktop/go-project/src/github.com/xxRanger/Percome19-Crowd-Demo/testnet/node/keystore/authority2"
const password = "321"


const myWallet = Wallet.fromV3(fs.readFileSync(utcFile).toString(), password, true);

console.log("Private Key: " + myWallet.getPrivateKey().toString('hex')) 
console.log("Address: " + myWallet.getAddress().toString('hex')) 