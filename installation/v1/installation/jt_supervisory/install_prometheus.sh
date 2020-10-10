#!/bin/bash
home="/jt/monitor"
path_install=$home"/install/jt_supervisory"
path_system="/lib/systemd/system"
name_prometheus="prometheus-2.17.2.linux-amd64"
path_prometheus=$home/$name_prometheus
file_name_prometheus_service="prometheus.service"
file_prometheus_sh=$path_prometheus"/start.sh"

#install prometheus
cd $home
pwd

#clear old versions
rm prometheus*.* -fr

#decompress
sudo tar -zxvf $path_install/$name_prometheus.tar.gz -C $home
cd $path_prometheus

#cp yml file
cp $path_install/prometheus.yml $path_prometheus/prometheus.yml

#create start script
echo './prometheus' >$file_prometheus_sh
chmod u+x $file_prometheus_sh

#create service file
(
cat <<EOF
[Unit]
Description=Prometheus Service
After=network-pre.target network-manager.service network-online.target network.target networking.service
[Service]
Type=idle
WorkingDirectory=$path_prometheus
User=root
ExecStart=/bin/bash $file_prometheus_sh
Restart=always
[Install]
WantedBy=multi-user.target

EOF
) >$path_system/$file_name_prometheus_service

#reload and run
sudo systemctl daemon-reload
sudo systemctl enable prometheus
sudo systemctl start prometheus
#sudo systemctl status prometheus
#sudo journalctl -f -u prometheus

cd $home


