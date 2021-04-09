const nano = require('nano-time');
const hashMap = require('hashmap');
const async = require('async');
let Web3 = require('web3');
let solc = require("solc");
let fs   = require('fs');
if(typeof web3 != 'undefined'){
	web3=new Web3(web3.currentProvider);
}else{
	web3 = new Web3('http://localhost:8545');
}
let source=fs.readFileSync("./demo.sol","utf8");
let cacl=solc.compile(source,1);
let abi= JSON.parse(cacl.contracts[':Calc'].interface);
let bytecode=cacl.contracts[':Calc'].bytecode;	
var txMap = new hashMap();
var txArray = new Array(500);




web3.eth.getAccounts().then(data=>{
	web3.eth.personal.unlockAccount(data[0],"1848252061",9999999).then(openAccountState=>{
		if(openAccountState){
			console.log("开户状态:"+openAccountState);
			var rsContract=new web3.eth.Contract(abi).deploy({
				data:'0x'+bytecode,
				arguments:[],	//传递构造函数的参数
			}).send({
				from:data[0],
				gas:1500000,
				gasPrice:'30000'
			},function(error,transactionHash){
				console.log("send回调");
				console.log("error:"+error);
				console.log("send transactionHash:"+transactionHash);
			})
			.on('error', function(error){ console.error(error) })
			// .on('transactionHash', function(transactionHash){ console.log("hash:",transactionHash)})
			// .on('receipt', function(receipt){
			//    console.log(receipt.contractAddress) // contains the new contract address
			// })
			//.on('confirmation', function(confirmationNumber, receipt){console.log("receipt,",receipt)})
			.then(function(newContractInstance){
				var newContractAddress=newContractInstance.options.address
				console.log("新合约地址:"+newContractAddress);
 
				web3.eth.getBlockNumber().then(blockNum=>{
					console.log("当前块号："+blockNum);
					web3.eth.getBlock(blockNum).then(data=>{
						console.log("当前块信息：");
						console.log(data);
					})
				});
				var MyContract = new web3.eth.Contract(abi,newContractAddress);
				function cal_add(sendTime,a,b){
					
					MyContract.methods.add(a,b).send({
						from:data[0],
						gas:15000000,
						gasPrice:'30000'
					},function(error,transactionHash){
						var respTime = nano();
						var time = {sendTime,respTime};
						txMap.set(transactionHash,time);
						txArray[a] = {transactionHash,sendTime,respTime};
						console.log("send transactionHash:"+transactionHash);
						console.log('tx size:'+txMap.size);
						console.log(txMap.get(transactionHash));
						if(txMap.size == 500){
						
						txMap.forEach(function(value,key){
							console.log(key+": "+value.sendTime);
						})
						function theReplacer(key, value) {
							if(key === "respTime"){
								return +value;
							}

							return key === "sendTime" ? +value : value;
						}

							fs.writeFile("./client2.json",JSON.stringify(txArray,theReplacer),function(err){
								if(err){
							  
								}
							  })
						}
					})
				}
				for (i = 0; i < 500; i++){
					cal_add(nano(),i,i+2);
				};
				// async function call_cal(){
				// 	return new Promise((resolve, reject) => {

				// 	resolve('cal');
				// 	});

				// }

				// async function collection(){
				// 	await call_cal();


				// 	txMap.forEach(function(value,key){
				// 		console.log(key+": "+value);
				// 	})
				// }
				// collection();

 
 
			});
			
		}
	});
});



