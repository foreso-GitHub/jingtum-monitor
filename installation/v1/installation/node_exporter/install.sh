#!/bin/bash
#check if current path is the path of sh
if [ ! -x install.sh ]; then
  echo -e "\033[41;36m Error: have to execute install.sh in its own path! \033[0m"
  exit 1
fi

home="/jt/monitor"
path_install=$home"/install/node_exporter"

#format in linux
sed -i 's/\r//' *.sh

#create install path and cp installations
rm $path_install -fr
mkdir -p $path_install
cp ./ $path_install -fr
cd $path_install
pwd

#start install
bash ./install_node.sh
