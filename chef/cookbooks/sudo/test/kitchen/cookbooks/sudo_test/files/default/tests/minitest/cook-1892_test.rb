#
# Cookbook Name:: sudo_test
# Minitest:: cook-1892
#
# Copyright 2012, Opscode, Inc.
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
#

require File.expand_path('../support/helpers', __FILE__)

describe "sudo_test::default" do
  include Helpers::SudoTest

  it 'creates a tomcat sudoers file' do
    file('/etc/sudoers.d/tomcat').must_exist
  end

  it 'has the correct permissions for tomcat' do
    if node['authorization']['sudo']['passwordless']
      file('/etc/sudoers.d/tomcat').must_include '%tomcat  ALL=(app_user) NOPASSWD:/etc/init.d/tomcat restart'
    else
      file('/etc/sudoers.d/tomcat').must_include '%tomcat  ALL=(app_user) /etc/init.d/tomcat restart'
    end
  end
end
