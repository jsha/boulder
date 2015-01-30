#!/bin/bash -ex
#sudo apt-get update
#wget https://storage.googleapis.com/golang/go1.4.1.linux-amd64.tar.gz
#tar xpvzf go1.4.1.linux-amd64.tar.gz
#cat <<EOF > ~/.bashrc
#export GOROOT=~/go GOPATH=~/go/packages
#EOF
#. ~/.bashrc

#CONF=/etc/acme
#sudo install -o vagrant -g vagrant -d $CONF
#openssl req -new -newkey rsa:2048 -nodes -days 10000 -x509 \
  #-keyout ${CONF}/acme.key -out ${CONF}/acme.crt \
  #-subj /CN=acme-dev-server
#cd /vagrant/
#go build boulder-start/
#./boulder-start monolithic
