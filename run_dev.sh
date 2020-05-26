#!/bin/bash

set -e

mkdir -p dev

cd server
go build -o ../dev/server
cd ..

cd dev
./server &
server_pid=$!

cd ../ui/
bash -c "NODE_ENV=development npm run serve" &
ui_pid=$!

function kill_servers {
    kill $server_pid
    kill $ui_pid
    wait
    exit
}

trap kill_servers SIGINT

wait
