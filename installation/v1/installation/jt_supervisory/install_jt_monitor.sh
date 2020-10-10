#!/bin/bash
home="/jt/monitor"
path_install=$home"/install/jt_supervisory"
path_system="/lib/systemd/system"
name_jt_monitor="jt_monitor-0.0.12.linux-amd64"
path_jt_monitor=$home/$name_jt_monitor
file_name_jt_monitor_service="jt_monitor.service"
file_jt_monitor_sh=$path_jt_monitor"/start.sh"

#install jt_monitor
cd $home
pwd

#clear old versions
rm jt_monitor*.* -fr

#decompress
sudo tar -zxvf $path_install/$name_jt_monitor.tar.gz
cd $path_jt_monitor

#create start script
echo './jt_monitor' >$path_jt_monitor/start.sh
chmod u+x $path_jt_monitor -R

#write service
(
cat <<EOF
[Unit]
Description=Jt Monitor Service
After=network-pre.target network-manager.service network-online.target network.target networking.service
[Service]
Type=idle
WorkingDirectory=$path_jt_monitor
User=root
ExecStart=/bin/bash $file_jt_monitor_sh
Restart=always
[Install]
WantedBy=multi-user.target

EOF
) >$path_system/$file_name_jt_monitor_service

#reload and run
sudo systemctl daemon-reload
sudo systemctl enable jt_monitor
sudo systemctl start jt_monitor
#sudo systemctl status jt_monitor
#sudo journalctl -f -u jt_monitor

cd $home

