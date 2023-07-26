#!/bin/bash
echo "[STEP1-1] UPDATE APT"
sudo su
apt-get update -y
apt install net-tools
ifconfig

#修改root账号密码
echo "[STEP1-2]Modify root password"
echo -e "zhongzn-1\nzhongzn-1" |sudo passwd root
sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
systemctl restart sshd

#修改hosts 
echo "[STEP2] SET HOSTS"
cat >>/etc/hosts<<EOF
192.168.56.10 k8s-master
192.168.56.11 k8s-node1
192.168.56.12 k8s-node2
EOF
cat /etc/hosts

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

#配置请求网桥
echo "[STEP4] IPTABLES SETTING"
cat >>/etc/sysctl.d/kubernetes.conf<<EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system >/dev/null 2>&1

#关闭内存置换
echo "[STEP5] DISABLE AND TURN OFF SWAP"
sed -i '/swap/d' /etc/fstab
swapoff -a
free -m

#解决搭建kubectl过慢问题
echo "[STEP6] INSTALL K8S RELAEASE"
apt-get update && apt-get install -y apt-transport-https curl
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add - 
cat >> /etc/apt/sources.list.d/kubernetes.list << EOF
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF
apt-get update
apt-get install -y kubelet kubeadm kubectl

#解决docker和k8s容器公用cgroup为systemed，可以通过docker info ｜grep cgroup查看
echo "[STEP7] MODIFY KUBELET SERVICE"
systemctl enable kubelet >/dev/null 2>&1
{
    "exec-opts": ["native.cgroupdriver=systemd"]
} >/dev/null 2>&1
systemctl restart kubelet

#解决containerd默认不适用cri
echo "[STEP8] CONTAINERD ENABLE"
sed -i 's/^disabled_plugins = \["cri"\]/#&/' /etc/containerd/config.toml
systemctl restart containerd

#解决证书过期问题
echo "[STEP9] CERTIFICATE EXPIRED"
apt install -y software-properties-common \
    gnupg debian-keyring debian-archive-keyring \
    apt-transport-https ca-certificates \
    lsb-core lsb-release
apt install -y ipset ipvsadm conntrack socat
apt-get update

# kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=192.168.56.10 --image-repository registry.aliyuncs.com/google_containers 

#安装go相关
VERSION=1.18 
wget https://studygolang.com/dl/golang/go${VERSION}.linux-amd64.tar.gz 
tar -xzvf go${VERSION}.linux-amd64.tar.gz 
sudo mv go /usr/local/go${VERSION} 
sudo ln -s /usr/local/go${VERSION}/bin/go /usr/bin/go
go env -w GOPROXY=https://goproxy.cn,direct
export GOPATH=`go env GOPATH`
sudo rm -f /etc/profile.d/99-go.sh
cat | sudo tee /etc/profile.d/99-go.sh << EOF
export GOPATH=$GOPATH
export PATH=\$PATH:\$GOPATH/bin
EOF
source /etc/profile
