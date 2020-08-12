#!/bin/bash
home="/jt/monitor"
#path_install="/opt"
path_system="/lib/systemd/system"
file_name_prometheus="prometheus-2.17.2.linux-amd64"
file_name_prometheus_service="prometheus.service"
url_prometheus="https://www.yiyuen.com/e/file/download?code=067c0871e41dac2b&id=27821"

#install prometheus
mkdir -p $home
cd $home
pwd

#clear old versions
rm prometheus*.* -fr

#download prometheus
output=`wget -O $file_name_prometheus.tar.gz $url_prometheus`
echo $output

#decompress
sudo tar -zxvf $home/$file_name_prometheus.tar.gz
cd $file_name_prometheus

#cp yml file
#cp $home/prometheus.yml $home/$file_name_prometheus/prometheus.yml

#create yml file
(
cat <<-EOF
# my global config
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
    - targets: ['localhost:9090']

  - job_name: "jt_node_al_1"
    static_configs:
    - targets: ["121.89.209.19:9100"]

  - job_name: "jt_node_al_2"
    static_configs:
    - targets: ["121.89.212.154:9100"]

  - job_name: "jt_node_al_3"
    static_configs:
    - targets: ["121.89.207.217:9100"]

  - job_name: "jt_node_al_4"
    static_configs:
    - targets: ["121.89.206.115:9100"]

  - job_name: "jt_node_al_5"
    static_configs:
    - targets: ["121.89.198.119:9100"]

  - job_name: "jt-node-bd"
    static_configs:
    - targets: ["180.76.125.22:9100"]

  - job_name: "jt-node-tx"
    static_configs:
    - targets: ["45.40.240.50:9100"]

  - job_name: "jt-node-hw"
    static_configs:
    - targets: ["121.37.216.100:9100"]

  - job_name: "jt-node-ty"
    static_configs:
    - targets: ["61.171.12.71:9100"]

EOF
) >$home/$file_name_prometheus/prometheus.yml

#create start script
echo './prometheus' >$home/$file_name_prometheus/start.sh
chmod u+x $home/$file_name_prometheus/start.sh

#create service file
(
cat <<EOF
[Unit]
Description=Prometheus Service
After=network-pre.target network-manager.service network-online.target network.target networking.service
[Service]
Type=idle
WorkingDirectory=$home/$file_name_prometheus/
User=root
ExecStart=/bin/bash $home/$file_name_prometheus/start.sh
Restart=always
[Install]
WantedBy=multi-user.target

EOF
) >$path_system/$file_name_prometheus_service

#reload and run
sudo systemctl daemon-reload
sudo systemctl enable prometheus
sudo systemctl start prometheus
sudo systemctl status prometheus
#sudo journalctl -f -u prometheus
cd $home


