const nano = require('nano-time');
const hashMap = require('hashmap');
const async = require('async');
var sleep = require('sleep');
let Web3 = require('web3');
let solc = require("solc");
let fs = require('fs');
if (typeof web3 != 'undefined') {
    web3 = new Web3(web3.currentProvider);
} else {
    web3 = new Web3('http://localhost:8545');
}
let source = fs.readFileSync("./sctest.sol", "utf8");
let cacl = solc.compile(source, 1);
let abi = JSON.parse(cacl.contracts[':kvstore'].interface);
let bytecode = cacl.contracts[':kvstore'].bytecode;
var txReqMap = new hashMap();
var txRespMap = new hashMap();
var txArray = new Array(500);


// var sleep = function(time) {
//     var startTime = new Date().getTime() + parseInt(time, 10);
//     while(new Date().getTime() < startTime) {}
// };

web3.eth.getAccounts().then(data => {
    web3.eth.personal.unlockAccount(data[0], "1848252061", 999999).then(openAccountState => {
        web3.eth.Contract.defaultAccount = web3.eth.accounts[0];
        console.log(web3.eth.Contract.defaultAccount)
        if (openAccountState) {
            var newContractAddress = '0x88446beb2169c0bbc7427874608cd85f2781c0fc'
            console.log("新合约地址:" + newContractAddress);

            web3.eth.getBlockNumber().then(blockNum => {
                console.log("当前块号：" + blockNum);
                web3.eth.getBlock(blockNum).then(data => {
                    // console.log("当前块信息：");
                    // console.log(data);
                })
            });
            var MyContract = new web3.eth.Contract(abi, newContractAddress);
            // function cal_add(sendTime,a,b){

            //     MyContract.methods.add(a,b).send({
            //         from:data[0],
            //         gas:15000000,
            //         gasPrice:'8100000000'
            //     },function(error,transactionHash){
            //         var respTime = nano();
            //         var time = {sendTime,respTime};
            //         txMap.set(transactionHash,time);
            //         txArray[a] = {transactionHash,sendTime,respTime};
            //         console.log("send transactionHash:"+transactionHash);
            //         console.log('tx size:'+txMap.size);
            //         console.log(txMap.get(transactionHash));

            //     })
            // }

            for (i = 1; i < 3; i++) {
                sleep.sleep(3);
                MyContract.methods.sort(i).send({
                    from: data[0],
                    gas: 3000000,
                    gasPrice: '3100000000'
                }, function (err, hash) {
                    console.log(hash);
                    txReqMap.set(hash, nano());

                }).then(function (receipt) {
                    // console.log(receipt);
                    console.log("receipt:" + receipt.transactionHash);
                    txRespMap.set(receipt.transactionHash, nano());
                    if (txRespMap.size == 5) {

                        txRespMap.forEach(function (value, key) {
                            console.log(key + ": " + value);
                        })
                        // function theReplacer(key, value) {
                        //     if(key === "respTime"){
                        //         return +value;
                        //     }

                        //     return key === "sendTime" ? +value : value;
                        // }
                        var options = { flag : 'w' };

                        fs.writeFile("./txReq.json", JSON.stringify(txReqMap), options,function (err) {
                            if (err) {

                            }
                        })
                        fs.writeFile("./txResp.json", JSON.stringify(txRespMap), options, function (err) {
                            if (err) {

                            }
                        })

                    }
                });

                // cal_add(nano(),i,i+2);
                // cal_add2(nano(),i,i+2);
            };


        }
    });
});



