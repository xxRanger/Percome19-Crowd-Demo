const Migrations = artifacts.require("Migrations");
const erc20 = artifacts.require("ERC20")
const crowd = artifacts.require("Crowdsourcing")
const fs = require('fs');
resolve = require('path').resolve
const backend_agg_path = resolve("../../back-end/etc/agg.json")
var agg_config = require(backend_agg_path);
const backend_erc20_path = resolve("../../back-end/etc/erc20.json")
var erc20_config = require(backend_erc20_path);

const frontend_agg_path = resolve("../../front-end/src/assets/config/agg.json")
var frontend_agg_config = require(frontend_agg_path);
var erc20_instance;
var agg_instance;


const account = {
                address: "0xcb4dF655d08C9f054493809DCDb5a47F3A95eFAA",
                privateKey: "7a318fc20fbfe6a35681ca1ba5d8a74a07e54af99c9942a38bb1b03aad6173c1"
            };

module.exports = function(deployer) {
  deployer.deploy(Migrations);
  deployer.deploy(erc20,10000000000,"crd","crd",0).then((instance)=>{
    erc20_instance = instance;
  	return deployer.deploy(crowd, erc20.address);
  }).then(async (instance)=>{
  		agg_instance = instance;
  		agg_config.address = crowd.address;
  		erc20_config.address = erc20.address;
  		frontend_agg_config.address = crowd.address;
  		 fs.writeFile(backend_agg_path,JSON.stringify(agg_config),(err)=>{if(err) console.log(err)});
  		 fs.writeFile(backend_erc20_path,JSON.stringify(erc20_config),(err)=>{if(err) console.log(err)});
  		 fs.writeFile(frontend_agg_path,JSON.stringify(frontend_agg_config),(err)=>{if(err) console.log(err)});
      let amount = 1000000;
      // console.log(erc20_instance);
      // let b = await erc20_instance.balanceOf("0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc")
      // console.log(b)
      await erc20_instance.transfer(account.address,amount,{from:"0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc"}).then(async (s)=>{
        // let b = await erc20_instance.balanceOf(account.address)
        // console.log("balance",b);
        // console.log(s);
        if(s.receipt.status!=true) throw "transfer fail";
        // console.log("test")
        return erc20_instance.approve(crowd.address, amount, {from: account.address})
      }).then(s=>{if(s.receipt.status!=true) throw "approve fail";})
  })
};

