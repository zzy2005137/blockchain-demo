#!/bin/bash

echo "组织1查询001号订单历史交易信息"
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["queryOrderHistoryPrivate","d4735e3a265e","d4735e3a265e","001"]}'

echo "组织2查询001号订单历史交易信息"
docker exec cli2 peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["queryOrderHistoryPrivate","d4735e3a265e","d4735e3a265e","001"]}'

