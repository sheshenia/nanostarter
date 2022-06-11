#!/bin/bash

mkdir ../scep && cd ../scep

curl -RLO https://github.com/micromdm/scep/releases/download/v2.1.0/scepserver-linux-amd64-v2.1.0.zip

unzip scepserver-linux-amd64-v2.1.0.zip

./scepserver-linux-amd64 ca -init