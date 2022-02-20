# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.hostname = "minitwit"
  config.vm.box = "digital_ocean"                                                 # create box
  config.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
  config.ssh.private_key_path = "~/.ssh/id_rsa"                                     # ssh private key location OBS '~/' is specific for Linux
  config.vm.synced_folder ".", "/vagrant", type: "rsync"                            # move all files

  
  #########################################################################
  #                                                                       #
  #                         MINITWIT BOX                                  #
  #                                                                       #
  #########################################################################
  config.vm.define "minitwit" do |server|
    server.vm.hostname = "minitwit"
    config.vm.provider :digital_ocean do |provider|
      provider.ssh_key_name = ENV["SSH_KEY_NAME"]
      provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
      provider.image = 'ubuntu-18-04-x64'
      provider.region = 'fra1'
      provider.size = 's-1vcpu-1gb'
      provider.privatenetworking = true
    end

    # provision go and setup envioremental variables
    # download -> unzip -> set PATH -> create new bash (to update variables)
    # tar -C <extraction path> -xzf (extract, gzip, file) <file>
    config.vm.provision "shell", inline: <<-SHELL
    wget https://go.dev/dl/go1.17.7.linux-amd64.tar.gz
      sudo tar -C /usr/local -xzf go1.17.7.linux-amd64.tar.gz
      echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
      source ~/.profile
    SHELL

    # provision gcc compiler for sqlite (gcc, g++ & make)
    config.vm.provision "shell", inline: <<-SHELL
      sudo apt -y install build-essential
    SHELL

    # provision sqlite3
    # download and copy original database to tmp
    config.vm.provision "shell", inline: <<-SHELL
      sudo apt -y install sqlite3
      cp /vagrant/og-db/minitwit.db /tmp/minitwit.db
    SHELL

    # download go dependicies from go.mod file
    # start the server
    config.vm.provision "shell", inline: <<-SHELL
      cd /vagrant
      go mod download
      go run src/minitwit.go
    SHELL
  end

  #########################################################################
  #                                                                       #
  #                           POSTGRES BOX                                #
  #                                                                       #
  #########################################################################
  #config.vm.define "postgresdb" do |server|
  #  config.vm.provider "virtualbox" do |provider|
  #    provider.gui = false
  #    provider.memory = "2048"
  #  end

    # provision for postgresql 
    # TODO: sudo apt -y install postgresql-12 and setup connection
  #  config.vm.provision "shell", inline: <<-SHELL
  #  SHELL
  #end

  # update
  config.vm.provision "shell", inline: <<-SHELL
    sudo apt-get update
  SHELL
end
