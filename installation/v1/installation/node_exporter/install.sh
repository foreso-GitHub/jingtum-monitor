#!/bin/bash
home="/jt/monitor"
path_install=$home"/install/node_exporter"

#format in linux
sed -i 's/\r//' install_node.sh

#create install path and cp installations
mkdir -p $path_install
cp ./ $path_install -fr
cd $path_install
pwd

#start install
bash ./install_node.sh
