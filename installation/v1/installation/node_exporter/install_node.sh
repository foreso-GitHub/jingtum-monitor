#!/bin/bash
home="/jt/monitor"
path_install=$home"/install/node_exporter"
path_system="/lib/systemd/system"
name_node_exporter="node_exporter-0.18.1.linux-amd64"
path_node_exporter=$home/$name_node_exporter
file_name_node_exporter_service="node_exporter.service"

#install node_exporter
cd $home
pwd

#clear old versions
rm node_exporter*.* -fr

#decompress
sudo tar -zxvf $path_install/$name_node_exporter.tar.gz -C $home
cd $path_node_exporter

#create start script
file_node_exporter_sh=$path_node_exporter/start.sh
echo './node_exporter' >$file_node_exporter_sh
chmod u+x $file_node_exporter_sh

#write service
(
cat <<EOF
[Unit]
Description=Node Exporter Service
After=network-pre.target network-manager.service network-online.target network.target networking.service
[Service]
Type=idle
WorkingDirectory=$path_node_exporter
User=root
ExecStart=/bin/bash $file_node_exporter_sh
Restart=always
[Install]
WantedBy=multi-user.target

EOF
) >$path_system/$file_name_node_exporter_service

#reload and run
sudo systemctl daemon-reload
sudo systemctl enable node_exporter
sudo systemctl start node_exporter
sudo systemctl status node_exporter
#sudo journalctl -f -u node_exporter

cd $home
