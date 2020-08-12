#!/bin/bash
name_install="node_exporter-amd64-install-0.1.0"
file_install=$name_install".tar.gz"
url="https://www.yiyuen.com/e/file/download?code=e93763121764211e&id=28380"

pwd
output=`wget -O $file_install $url`
echo $output

sudo tar -zxvf $file_install
cd $name_install
sed -i 's/\r//' install.sh
bash ./install.sh
