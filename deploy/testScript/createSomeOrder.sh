#!/bin/bash

#create order
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["createOrder","5feceb66ffc8","d4735e3a265e","003","toBeStarted"]}'  #管理员id, owner, orderId, status
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["createOrder","5feceb66ffc8","d4735e3a265e","004","inProgress"]}'  #管理员id, owner, orderId, status
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["createOrder","5feceb66ffc8","d4735e3a265e","005","toBeStarted"]}'  #管理员id, owner, orderId, status
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["createOrder","5feceb66ffc8","d4735e3a265e","006","done"]}'  #管理员id, owner, orderId, status
