#!/bin/bash
home="/jt/monitor"
path_install=$home"/install/jt_supervisory"

#format in linux
sed -i 's/\r//' install_all.sh

#create install path and cp installations
rm $path_install -fr
mkdir -p $path_install

#todo if current path is the path of sh
cp ./ $path_install -fr
cd $path_install
pwd

#start install
bash ./install_all.sh
