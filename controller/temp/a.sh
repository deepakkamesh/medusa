#!/bin/bash

case $1 in
 init)
   export IP="localhost"
   export PADDR="68,65,6C,6C,6F"
   export ADDR="FF,FF,FF"
   ;;
 prt)
   echo $PADDR
   echo $ADDR
   ;;
 env)
   export PADDR=$2
   export ADDR=$3
   ;;
 test)
   echo "test func"
   go run ../cli.go -- $IP:8080 T AD,$PADDR,10,$ADDR,16;
   ;;
 volt)
   echo "voltage"
   go run ../cli.go -- $IP:8080 T AD,$PADDR,10,$ADDR,5;
   ;;
 lon)
   echo "turn on led"
   go run ../cli.go -- $IP:8080 T AD,$PADDR,10,$ADDR,13,1;
   ;;
 lof)
   echo "turn off led"
   go run ../cli.go -- $IP:8080 T AD,$PADDR,10,$ADDR,13,0;
   ;;
 rlon)
   echo "turn on board led"
   go run ../cli.go -- $IP:8080 T AD,0,0,0,0,0,10,$ADDR,13,0;
   ;;
 rlof)
   echo "turn off led"
   go run ../cli.go -- $IP:8080 T AD,0,0,0,0,0,10,$ADDR,13,1;
   ;;
 flush)
   echo "flushing relay tx"
   go run ../cli.go -- $IP:8080 T AD,0,0,0,0,0,10,$ADDR,17;
   ;;
 rsb)
   echo "restarting board"
   go run ../cli.go -- $IP:8080 T AD,$PADDR,10,$ADDR,14;
   ;;
 rsr)
   echo "restarting relay"
   go run ../cli.go -- $IP:8080 T AD,0,0,0,0,0,10,$ADDR,14;
   ;;
 stress)
   echo "stress test board"
   for i in {1..11}; do go run ../cli.go -- $IP:8080 T AD,$PADDR,10,$ADDR,21;done
   ;;
 bc)
   # $2 - Pipe Addr, $3 - Board Addr
   echo "change board config"
   go run ../cli.go -- $IP:8080 T AD,$PADDR,5,$ADDR,A,2,$2,$3
   ;;
 rc)
   echo "sending relay cfg"
   go run ../cli.go -- $IP:8080 T AB,$2,A,1,2,3,4,B,C,D,E,74,A,B,C
   ;;
 rc1)
   echo "sending relay cfg"
   go run ../cli.go -- $IP:8080 U AB,68,65,6C,6C,6F,A,1,2,3,4,B,C,D,E,74,A,B,C
   ;;
 relaycfg1)
esac
