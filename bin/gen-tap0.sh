#!/bin/sh

set -e

ip tuntap add mode tap user $USER name tap0
ip addr add 192.0.2.1/16 dev tap0
ip link set tap0 up
