offline
解压安装包，在目录中执行install.sh

online
sh下载offline安装包


script清单：
1. install_node.sh：copy这个文件到节点服务器后执行（或者copy其中的script直接在节点服务器直接执行）。script会自动下载安装node_exporter，并用systemd自动运行。
2. node_exporter.service：node_exporter对应systemd的配置文件，由install_node.sh自动生成。
3. install_promethues.sh：
4. 


步骤：

1. 安装node_exporter
1.1 打开所有节点的9100端口
1.2 更新node_exporter.service
1.3 更新install_node_exporter.sh
1.4 每个节点执行install_node_exporter.sh
1.5 node_exporter运行成功

2. 安装prometheus
2.1 打开节点的9090端口
2.2 更新prometheus.yml
2.3 更新prometheus.service
2.4 更新install_promethues.sh有效链接
2.5 执行install_promethues.sh
2.6 promethues运行成功

3. 安装grafana
3.1 打开节点的3000端口
3.2 更新grafana.service
3.3 更新install_grafana.sh有效链接
3.4 执行install_grafana.sh
3.5 grafana运行成功

4. 安装jt_monitor
1.1 打开节点的9101端口


