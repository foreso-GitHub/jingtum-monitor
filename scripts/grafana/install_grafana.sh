#!/bin/bash
home="/jt/monitor"


#path_install="/opt"
path_system="/lib/systemd/system"
file_name_grafana="grafana-6.7.3"
file_name_grafana_service="grafana.service"
url_grafana="https://www.yiyuen.com/e/file/download?code=8ead1816bbc43aa5&id=27827"

#install grafana
mkdir -p $home
cd $home
pwd
#clear old versions
rm grafana*.* -fr
#download grafana
output=`wget -O $file_name_grafana.tar.gz $url_grafana`
echo $output
#decompress
sudo tar -zxvf $home/$file_name_grafana.tar.gz
cd $file_name_grafana
#create start script
echo './bin/grafana-server web' >$home/$file_name_grafana/start.sh
chmod u+x $home/$file_name_grafana/start.sh
#create service file
(
cat <<EOF
[Unit]
Description=grafana Service
After=network-pre.target network-manager.service network-online.target network.target networking.service
[Service]
Type=idle
WorkingDirectory=$home/$file_name_grafana/
User=root
ExecStart=/bin/bash $home/$file_name_grafana/start.sh
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
