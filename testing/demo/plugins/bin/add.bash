#!/usr/bin/env bash
num=$1
shift
delta=$1

ret=`expr $num + $delta`
echo "$ret"