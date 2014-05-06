#!/usr/bin/ruby
# encoding: utf-8

# NOTE: proper serialization of yaml literal blocks seems to require ruby 2.0
require 'yaml'

def prepare_units
  units = []
  source_dirs = [
    File.join('contrib/coreos/systemd'),
    #File.join('controller/systemd'),
    #File.join('cache/systemd'),
    #File.join('database/systemd'),
    #File.join('logger/systemd'),
    #File.join('registry/systemd'),
    #File.join('router/systemd'),
  ]
  source_dirs.each do |dir|
    Dir.glob(File.join(dir, '*')) do |path|
      name = File.basename(path)
      unit = {'name' => name, 'command' => 'start',
              'content' => IO.read(path)}
      units.push(unit)
    end
  end
  return units
end

# join generated/default user-data and write to stdout
user_data = YAML.load_file('contrib/coreos/default-user-data.yml')
user_data['coreos']['units'].concat(prepare_units())
puts "#cloud-config\n"
puts user_data.to_yaml
