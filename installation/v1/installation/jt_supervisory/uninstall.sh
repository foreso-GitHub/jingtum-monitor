#!/bin/bash
home="/jt/monitor"
path_install=$home"/install/jt_supervisory"
path_system="/lib/systemd/system"

name_jt_monitor="jt_monitor-0.0.7.linux-amd64"
path_jt_monitor=$home/$name_jt_monitor
file_name_jt_monitor_service="jt_monitor.service"

name_prometheus="prometheus-2.17.2.linux-amd64"
path_prometheus=$home/$name_prometheus
file_name_prometheus_service="prometheus.service"

name_grafana="grafana-6.7.3"
path_grafana=$home/$name_grafana
file_name_grafana_service="grafana.service"

sudo systemctl stop grafana
rm $path_grafana -fr
rm $path_system/$file_name_grafana_service

sudo systemctl stop prometheus
rm $path_prometheus -fr
rm $path_system/$file_name_prometheus_service

sudo systemctl stop jt_monitor
rm $path_jt_monitor -fr
rm $path_system/$file_name_jt_monitor_service

rm $path_install -fr
