# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure("2") do |config|
  # Host and guest must be the same architecture, since boulder-start is built
  # on the host and distributed to all guests.
  config.vm.box = "trusty64"
  config.vm.box_url = "https://cloud-images.ubuntu.com/vagrant/trusty/20150123/trusty-server-cloudimg-amd64-vagrant-disk1.box"
  config.vm.box_download_checksum_type = "sha256"
  config.vm.box_download_checksum = "2d3f4bd2b3b9e84fb3eacf5034c778dc89c41505f010469126674124bf3ad880"
  config.vm.provision :shell, privileged: false, path: "vagrant-up.sh"

  config.vm.synced_folder ENV["GOPATH"].split(":")[0] + "/bin", "/home/vagrant/bin"

  config.vm.provider :virtualbox do |vb|
    vb.customize ["modifyvm", :id, "--memory", 256]
  end

  config.vm.define "frontend" do |frontend|
    frontend.vm.network :forwarded_port, host: 4000, guest: 4000
  end
  config.vm.define "va" do |va|
  end
end
