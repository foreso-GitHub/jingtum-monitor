#!/bin/bash
sed -i 's/\r//' install_jt_monitor.sh && \
sed -i 's/\r//' install_prometheus.sh && \
sed -i 's/\r//' install_grafana.sh && \
bash ./install_jt_monitor.sh && bash ./install_prometheus.sh && bash ./install_grafana.sh