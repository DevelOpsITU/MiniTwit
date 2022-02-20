# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/trusty64"                                                 # create box
  config.vm.network "forwarded_port", guest: 8080, host: 8080, host_ip: "127.0.0.1" # setup ips
  config.vm.synced_folder ".", "/vagrant", type: "rsync"                            # move all files

  config.vm.provider "virtualbox" do |server|
    server.gui = false
    server.memory = "2048"
  end

  # provision go and setup envioremental variables
  # wget download zip
  # tar -C <extraction path> -xzf (extract, gzip, file) <file>
  config.vm.provision "shell", inline: <<-SHELL
    wget "go1.17.7.linux-amd64.tar.gz"
    sudo tar -C /usr/local/ -xzf go1.17.7.linux-amd64.tar.gz
    echo "
    export GOPATH=$HOME/go
    export GOROOT=/usr/local/go
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    "
  SHELL

  # start the server
  config.vm.provision "shell", inline: <<-SHELL
    go run src/minitwit.go
  SHELL
end
