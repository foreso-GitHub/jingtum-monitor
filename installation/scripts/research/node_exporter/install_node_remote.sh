#!/bin/bash
home="/jt/monitor"
mkdir $home
cd $home

url_node_exporter_install="https://www.yiyuen.com/e/file/download?code=c6635ef1a15627f1&id=27742"
file_name_node_exporter_install="install_node.sh"
output=`wget -O $file_name_node_exporter_install $url_node_exporter_install`
echo $output
chmod u+x $home/$file_name_node_exporter_install
sed -i 's/\r//' $home/$file_name_node_exporter_install
bash $home/$file_name_node_exporter_install




