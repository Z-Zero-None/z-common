#!/bin/bash
curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh" | sudo bash
sudo apt-get install gitlab-runner
#使用注册的gitlab-runner：gitlab-runner register  --url http://192.168.56.120  --token glrt--3fh59JXa1sWRWMinQxk
