#!/bin/bash
home="/jt/monitor"

mkdir -p $home
cd $home
pwd
output=`wget -O install_node.sh https://www.yiyuen.com/e/file/download?code=1255bc51c25b566c&id=27886`
echo $output
output=`wget -O install_prometheus.sh https://www.yiyuen.com/e/file/download?code=2ae36f937de343d6&id=27885`
echo $output
output=`wget -O install_grafana.sh https://www.yiyuen.com/e/file/download?code=d6c737e285e002e4&id=27884 `
echo $output
sed -i 's/\r//' install_node.sh && \
sed -i 's/\r//' install_prometheus.sh && \
sed -i 's/\r//' install_grafana.sh && \
bash ./install_node.sh && bash ./install_prometheus.sh && bash ./install_grafana.sh



