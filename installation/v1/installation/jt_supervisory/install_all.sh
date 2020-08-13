#!/bin/bash
sed -i 's/\r//' *.sh
bash ./install_jt_monitor.sh && bash ./install_prometheus.sh && bash ./install_grafana.sh