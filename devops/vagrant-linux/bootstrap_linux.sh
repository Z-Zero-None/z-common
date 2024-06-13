#!/bin/bash
echo "[STEP1-1] UPDATE APT"
sudo su
apt-get update -y
apt install net-tools
ifconfig
apt-get install ntpdate
ntpdate ntp1.aliyun.com
#修改root账号密码
echo "[STEP1-2]Modify root password"
echo -e "zhongzn-1\nzhongzn-1" |sudo passwd root
sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
systemctl restart sshd
#安装符合当前Linux操作系统的docker引擎以及相关证书
echo "[STEP3-1] INSTALL DOCKER"
apt-get install apt-transport-https ca-certificates curl software-properties-common -y
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
apt-get update -y
apt-get install docker-ce -y
usermod -aG docker $USER
#启动docker
echo "[STEP3-2]MODIFY DOCKER IMAGE-SOURCE"
systemctl enable docker >/dev/null 2>&1
systemctl start docker
#解决镜像拉去过慢
echo "[STEP3-3] ENABLE AND START DOCKER"
cat >>/etc/docker/daemon.json<<EOF
{
    "log-level":        "error" ,
    "exec-opts": ["native.cgroupdriver=systemd"],
    "registry-mirrors": [
        "https://mirror.ccs.tencentyun.com",
        "https://registry.docker-cn.com",
        "https://hub-mirror.c.163.com"
      ]
}
EOF
systemctl restart docker
docker -v

#安装gitlab
export GITLAB_HOME=/srv/gitlab
docker compose -d -f /srv/docker-gitlab/docker-compose.yml up
