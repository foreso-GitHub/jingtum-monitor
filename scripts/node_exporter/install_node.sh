#!/bin/bash
home="/jt/monitor"
mkdir -p $home
cd $home

path_install="/opt"
path_system="/lib/systemd/system"
file_name_node_exporter="node_exporter-0.18.1.linux-amd64"
file_name_node_exporter_service="node_exporter.service"
#url_node_exporter="https://www.yiyuen.com/e/file/download?code=6b2be671e7aadc4c&id=27692"
url_node_exporter="https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz"


#install node_exporter
cd $path_install
pwd
#clear old versions
rm node_exporter*.* -fr
#download node exporter
output=`wget -O $file_name_node_exporter.tar.gz $url_node_exporter`
echo $output
#decompress
sudo tar -zxvf $path_install/$file_name_node_exporter.tar.gz
cd $file_name_node_exporter
#create start script
echo './node_exporter' >$path_install/$file_name_node_exporter/start.sh
chmod u+x $path_install/$file_name_node_exporter/start.sh
#write service
(
cat <<EOF
[Unit]
Description=Node Exporter Service
After=network-pre.target network-manager.service network-online.target network.target networking.service
[Service]
Type=idle
WorkingDirectory=$path_install/$file_name_node_exporter/
User=root
ExecStart=/bin/bash $path_install/$file_name_node_exporter/start.sh
Restart=always
[Install]
WantedBy=multi-user.target

EOF
) >$path_system/$file_name_node_exporter_service
#reload and run
sudo systemctl daemon-reload
sudo systemctl start node_exporter
sudo systemctl status node_exporter
sudo journalctl -f -u node_exporter


