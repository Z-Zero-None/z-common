#!/bin/bash
# https://about.gitlab.com/install/#ubuntu
sudo apt-get update
sudo apt-get install -y curl openssh-server ca-certificates tzdata perl
sudo apt-get install -y postfix
curl https://packages.gitlab.com/install/repositories/gitlab/gitlab-ee/script.deb.sh | sudo bash
sudo EXTERNAL_URL="http://192.168.56.120" apt-get install gitlab-ee
# 可以通过 cat /etc/gitlab/gitlab.rb
# 对external_url "http://10.0.0.1"进行调整
# 重启配置 sudo gitlab-ctl reconfigure
# cat /etc/gitlab/initial_root_password 修改root密码为zhongzn-1
# 创建项目生成token：glpat-dNatzrtzzMF5TGTboHdX
# 创建gitlab-runner生成runner token

#启动gitlab内部的docker容器
#修改/etc/gitlab/gitlab.rb配置 registry_external_url 'https://192.168.56.120:5050'
#重启配置 gitlab-ctl reconfigure
#查看配置 gitlab-ctl status
#test:openssl s_client -showcerts -servername 192.168.56.120 -connect 192.168.56.120:5050 > cacert.pem