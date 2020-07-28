#!/bin/bash
home="/jt/monitor"
mkdir $home
cd $home

url_node_exporter_install="https://www.yiyuen.com/e/file/download?code=c6635ef1a15627f1&id=27742"
file_name_node_exporter_install="install_node.sh"

(
cat <<EOF
#!/bin/bash
mkdir $home
cd $home
echo `wget -O $file_name_node_exporter_install $url_node_exporter_install`
chmod u+x $home/$file_name_node_exporter_install
sed -i 's/\r//' $home/$file_name_node_exporter_install
bash $home/$file_name_node_exporter_install

EOF
) >$home/remote.sh

chmod u+x $home/remote.sh
bash $home/remote.sh

# output='scp -p $home/remote.sh root@121.89.198.119:/root/remote.sh'
# echo $output
# output='ssh -t root@121.89.198.119 sudo /root/remote.sh'
# echo $output


