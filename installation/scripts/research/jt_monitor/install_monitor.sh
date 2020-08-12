#!/bin/bash
home="/jt/monitor"
# path_install="/opt"
path_system="/lib/systemd/system"
file_name_jt_monitor="jt_monitor-0.0.7.linux-amd64"
file_name_jt_monitor_service="jt_monitor.service"
url_jt_monitor="https://www.yiyuen.com/e/file/download?code=cc958783301c62b2&id=28053"


#install jt_monitor
mkdir -p $home
cd $home
pwd
#clear old versions
rm jt_monitor*.* -fr
#download node exporter
output=`wget -O $file_name_jt_monitor.tar.gz $url_jt_monitor`
echo $output
#decompress
sudo tar -zxvf $home/$file_name_jt_monitor.tar.gz
cd $file_name_jt_monitor
#create start script
echo './jt_monitor' >$home/$file_name_jt_monitor/start.sh
chmod u+x $home/$file_name_jt_monitor -R
#chmod u+x $home/$file_name_jt_monitor/start.sh
#write service
(
cat <<EOF
[Unit]
Description=Jt Monitor Service
After=network-pre.target network-manager.service network-online.target network.target networking.service
[Service]
Type=idle
WorkingDirectory=$home/$file_name_jt_monitor/
User=root
ExecStart=/bin/bash $home/$file_name_jt_monitor/start.sh
Restart=always
[Install]
WantedBy=multi-user.target

EOF
) >$path_system/$file_name_jt_monitor_service
#reload and run
sudo systemctl daemon-reload
sudo systemctl enable jt_monitor
sudo systemctl start jt_monitor
sudo systemctl status jt_monitor
#sudo journalctl -f -u jt_monitor

cd $home

