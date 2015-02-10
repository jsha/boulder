# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

if !ENV["GOPATH"]
  print "Must set $GOPATH first."
  exit 1
end

GO_BIN = ENV["GOPATH"].split(":")[0] + "/bin"
if !File.exists?(GO_BIN + "/boulder-start")
  print "Must run go install github.com/letsencrypt/boulder/boulder-start"
  exit 1
end

Vagrant.configure("2") do |config|
  # Host and guest must be the same architecture, since boulder-start is built
  # on the host and distributed to all guests.
  config.vm.box = "trusty64"
  config.vm.box_url = "https://cloud-images.ubuntu.com/vagrant/trusty/20150123/trusty-server-cloudimg-amd64-vagrant-disk1.box"
  config.vm.box_download_checksum_type = "sha256"
  config.vm.box_download_checksum = "2d3f4bd2b3b9e84fb3eacf5034c778dc89c41505f010469126674124bf3ad880"

  config.vm.synced_folder ".", "/vagrant"
  config.vm.synced_folder GO_BIN, "/home/vagrant/bin"

  config.vm.provider :virtualbox do |vb|
    vb.customize ["modifyvm", :id, "--memory", 256]
  end

  config.vm.define "wfe" do |wfe|
    wfe.vm.network :forwarded_port, host: 4000, guest: 4000
    wfe.vm.provision :shell, privileged: false, inline: "~/bin/boulder-start --listen=0.0.0.0:4000 wfe"
  end
  config.vm.define "va" do |va|
    va.vm.provision :shell, privileged: false, inline: "~/bin/boulder-start va"
  end
  config.vm.define "ca" do |ca|
    ca.vm.provision :shell, privileged: false, inline: "~/bin/boulder-start ca"
  end
  config.vm.define "sa" do |sa|
    sa.vm.provision :shell, privileged: false, inline: "~/bin/boulder-start sa"
  end
  config.vm.define "ra" do |ra|
    ra.vm.provision :shell, privileged: false, inline: "~/bin/boulder-start ra"
  end
  config.vm.define "rabbitmq" do |rabbitmq|
    rabbitmq.vm.provision :shell, privileged: false, path: "rabbitmq/rabbitmq-setup.sh"
  end
end
