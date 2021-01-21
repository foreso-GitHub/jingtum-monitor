#!/bin/bash
name_install="jt_supervisory-amd64-install-0.1.4"
file_install=$name_install".tar.gz"
url="http://60.179.32.90:2811/e/file/download?code=d15d9397314caaaa&id=45445"

pwd
output=`wget -O $file_install $url`
echo $output

sudo tar -zxvf $file_install
cd $name_install
sed -i 's/\r//' install.sh
bash ./install.sh
