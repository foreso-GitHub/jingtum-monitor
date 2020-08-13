#!/bin/bash
home="/jt/monitor"
path_install=$home"/install/jt_supervisory"
path_system="/lib/systemd/system"
name_grafana="grafana-6.7.3"
file_name_grafana="grafana-6.7.3.linux-amd64"
path_grafana=$home/$name_grafana
file_name_grafana_service="grafana.service"
file_grafana_sh=$path_grafana"/start.sh"

#install grafana
cd $home
pwd

#clear old versions
rm grafana*.* -fr

#decompress
sudo tar -zxvf $path_install/$file_name_grafana.tar.gz -C $home
cd $path_grafana

#create start script
echo './bin/grafana-server web' >$file_grafana_sh
chmod u+x $file_grafana_sh

#create service file
(
cat <<EOF
[Unit]
Description=grafana Service
After=network-pre.target network-manager.service network-online.target network.target networking.service
[Service]
Type=idle
WorkingDirectory=$path_grafana
User=root
ExecStart=/bin/bash $file_grafana_sh
Restart=always
[Install]
WantedBy=multi-user.target

EOF
) >$path_system/$file_name_grafana_service

#reload and run
sudo systemctl daemon-reload
sudo systemctl enable grafana
sudo systemctl start grafana
sudo systemctl status grafana
#sudo journalctl -f -u grafana

cd $home
