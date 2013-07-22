#
# Author:: Bryan W. Berry (<bryan.berry@gmail.com>)
# Cookbook Name:: sudo
# Resource:: default
#
# Copyright 2011, Bryan w. Berry
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

actions :install, :remove
default_action :install

attribute :user,       :kind_of => String,          :default => nil
attribute :group,      :kind_of => String,          :default => nil
attribute :commands,   :kind_of => Array,           :default => ['ALL']
attribute :host,       :kind_of => String,          :default => 'ALL'
attribute :runas,      :kind_of => String,          :default => 'ALL'
attribute :nopasswd,   :equal_to => [true, false],  :default => false
attribute :template,   :regex => /^[a-z_]+.erb$/,   :default => nil
attribute :variables,  :kind_of => Hash,            :default => nil

# Set default for the supports attribute in initializer since it is
# a 'reserved' attribute name
def initialize(*args)
  super
  @action = :install
  @supports = { :report => true, :exception => true }
end
