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

#修改hosts 
echo "[STEP2] SET HOSTS"
cat >>/etc/hosts<<EOF
192.168.56.10 k8s-master
192.168.56.11 k8s-node1
192.168.56.12 k8s-node2
192.30.253.119 gist.github.com
54.169.195.247 api.github.com
185.199.111.153 assets-cdn.github.com
151.101.64.133 raw.githubusercontent.com
151.101.108.133 user-images.githubusercontent.com
151.101.76.133 gist.githubusercontent.com
151.101.76.133 cloud.githubusercontent.com
151.101.76.133 camo.githubusercontent.comkub
151.101.76.133 avatars0.githubusercontent.com
151.101.76.133 avatars1.githubusercontent.com
151.101.76.133 avatars2.githubusercontent.com
151.101.76.133 avatars3.githubusercontent.com
151.101.76.133 avatars4.githubusercontent.com
151.101.76.133 avatars5.githubusercontent.com
151.101.76.133 avatars6.githubusercontent.com
151.101.76.133 avatars7.githubusercontent.com
151.101.76.133 avatars8.githubusercontent.com
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

#wget https://github.com/containerd/containerd/releases/download/v1.7.0/containerd-1.7.0-linux-amd64.tar.gz
#tar Cxzvf /usr/local containerd-1.7.0-linux-amd64.tar.gz
cat << EOF >> /lib/systemd/system/containerd.service
[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target local-fs.target
​
[Service]
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/containerd
​
Type=notify
Delegate=yes
KillMode=process
Restart=always
RestartSec=5
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNPROC=infinity
LimitCORE=infinity
LimitNOFILE=infinity
# Comment TasksMax if your systemd version does not supports it.
# Only systemd 226 and above support this version.
TasksMax=infinity
OOMScoreAdjust=-999
​
[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable --now containerd
#配置请求网桥
echo "[STEP4] IPTABLES SETTING"
modprobe overlay
modprobe br_netfilter
cat <<EOF >> /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

cat >>/etc/sysctl.d/kubernetes.conf<<EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward=1
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

#解决containerd默认不适用cri+替换镜像选用
echo "[STEP8] CONTAINERD ENABLE"
containerd config default | tee /etc/containerd/config.toml
sed -i 's/registry.k8s.io\/pause:3.6/registry.aliyuncs.com\/google_containers\/pause:3.6/g' /etc/containerd/config.toml
sed -i "s#k8s.gcr.io#registry.cn-hangzhou.aliyuncs.com/google_containers#g"  /etc/containerd/config.toml
#sed -i "s#https://registry-1.docker.io#https://registry.cn-hangzhou.aliyuncs.com#g"  /etc/containerd/config.toml
#sed -i 's/^disabled_plugins = \["cri"\]/#&/' /etc/containerd/config.toml
cat /etc/containerd/config.toml | grep -n "sandbox_image"
systemctl restart containerd && systemctl status containerd

#解决crictl images 出现报错问题"ListImages with filter from image service failed"
crictl config runtime-endpoint unix:///run/containerd/containerd.sock
crictl config image-endpoint unix:///run/containerd/containerd.sock

#解决证书过期问题
echo "[STEP9] CERTIFICATE EXPIRED"
apt install -y software-properties-common \
    gnupg debian-keyring debian-archive-keyring \
    apt-transport-https ca-certificates \
    lsb-core lsb-release
apt install -y ipset ipvsadm conntrack socat
apt-get update

# kubeadm init --pod-network-cidr=10.244.0.0/16 --kubernetes-version=v1.27.3 --apiserver-advertise-address=192.168.56.10 --image-repository registry.aliyuncs.com/google_containers

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

#echo -e "y" |sudo ufw enable
#ufw allow 6443/tcp

#处理container中pause：3.6问题
#containerd config default > /etc/containerd/config.toml
#cat /etc/containerd/config.toml | grep -n "sandbox_image"
##替换镜像 sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.6"
#sed -i 's/registry.k8s.io\/pause:3.6/registry.aliyuncs.com\/google_containers\/pause:3.6/g' /etc/containerd/config.toml
#cat /etc/containerd/config.toml | grep -n "sandbox_image"
systemctl daemon-reload
systemctl restart containerd

