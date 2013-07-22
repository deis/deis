require 'spec_helper'

describe 'sudo::default' do
  context 'usual business' do
    before { Fauxhai.mock :platform => 'ubuntu' }
    let(:runner) { ChefSpec::ChefRunner.new.converge 'sudo::default' }

    it 'installs the sudo package' do
      runner.should install_package 'sudo'
    end

    it 'creates the /etc/sudoers file' do
      runner.should create_file_with_content '/etc/sudoers', 'Defaults      !lecture,tty_tickets,!fqdn'
    end
  end

  context 'sudoers.d' do
    let(:runner) do
      ChefSpec::ChefRunner.new do |node|
        node.set['authorization'] = {
          'sudo' => {
            'include_sudoers_d' => 'true'
          }
        }
      end.converge 'sudo::default'
    end

    it 'creates the sudoers.d directory' do
      runner.should create_directory '/etc/sudoers.d'
      runner.directory('/etc/sudoers.d').should be_owned_by 'root', 'root'
    end

    it 'drops the README file' do
      runner.should create_file_with_content '/etc/sudoers.d/README', 'As of Debian version 1.7.2p1-1, the default /etc/sudoers file created on'
    end
  end
end
