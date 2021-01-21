#!/bin/bash
name_install="jt_supervisory-amd64-install-0.1.4"
file_install=$name_install".tar.gz"
url="http://60.179.35.215:2811/e/file/download?code=b1dcc751639e1ba8&id=45445"

pwd
output=`wget -O $file_install $url`
echo $output

sudo tar -zxvf $file_install
cd $name_install
sed -i 's/\r//' install.sh
bash ./install.sh
