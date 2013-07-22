#
# Author:: Bryan W. Berry (<bryan.berry@gmail.com>)
# Author:: Seth Vargo (<sethvargo@gmail.com>)
# Cookbook Name:: sudo
# Provider:: default
#
# Copyright 2011, Bryan w. Berry
# Copyright 2012, Seth Vargo
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

# This LWRP supports whyrun mode
def whyrun_supported?
  true
end

# Ensure that the inputs are valid (we cannot just use the resource for this)
def check_inputs(user, group, foreign_template, foreign_vars)
  # if group, user, and template are nil, throw an exception
  if user == nil && group == nil && foreign_template == nil
    Chef::Application.fatal!('You must provide a user, group, or template!')
  elsif user != nil && group != nil && template != nil
    Chef::Application.fatal!('You cannot specify user, group, and template!')
  end
end

# Validate the given resource (template) by writing it out to a file and then
# ensuring that file's contents pass `visudo -c`
def validate_fragment!(resource)
  file = Tempfile.new('sudoer')

  begin
    file.write(capture(resource))

    cmd = Mixlib::ShellOut.new("visudo -cf #{file.path}").run_command
    unless cmd.exitstatus == 0
      Chef::Log.error("Fragment validation failed: \n\n")
      Chef::Log.error(file.read)
      Chef::Application.fatal!("Template #{file.path} failed fragment validation!")
    end
  ensure
    file.close
    file.unlink
  end
end

# Render a single sudoer template. This method has two modes:
#   1. using the :template option - the user can specify a template
#      that exists in the local cookbook for writing out the attributes
#   2. using the built-in template (recommended) - simply pass the
#      desired variables to the method and the correct template will be
#      written out for the user
def render_sudoer
  if new_resource.template
    Chef::Log.debug('Template attribute provided, all other attributes ignored.')

    resource = template "/etc/sudoers.d/#{new_resource.name}" do
      source        new_resource.template
      owner         'root'
      group         'root'
      mode          '0440'
      variables     new_resource.variables
      action        :nothing
    end
  else
    sudoer = new_resource.user || "%#{new_resource.group}".squeeze('%')

    resource = template "/etc/sudoers.d/#{new_resource.name}" do
      source        'sudoer.erb'
      cookbook      'sudo'
      owner         'root'
      group         'root'
      mode          '0440'
      variables     :sudoer => sudoer,
                    :host => new_resource.host,
                    :runas => new_resource.runas,
                    :nopasswd => new_resource.nopasswd,
                    :commands => new_resource.commands
      action        :nothing
    end
  end

  # Ensure that, adding this sudoer, would not break sudo
  validate_fragment!(resource)

  resource.run_action(:create)
  new_resource.updated_by_last_action(true) if resource.updated_by_last_action?
end

# Default action - install a single sudoer
action :install do
  render_sudoer
end

# Removes a user from the sudoers group
action :remove do
  resource = file "/etc/sudoers.d/#{new_resource.name}" do
    action :nothing
  end
  resource.run_action(:delete)
  new_resource.updated_by_last_action(true) if resource.updated_by_last_action?
end

private
# Capture a template to a string
def capture(template)
  context = {}
  context.merge!(template.variables)
  context[:node] = node

  eruby = Erubis::Eruby.new(::File.read(template_location(template)))
  return eruby.evaluate(context)
end

# Find the template
def template_location(template)
  if template.local
    template.source
  else
    context = template.instance_variable_get('@run_context')
    cookbook = context.cookbook_collection[template.cookbook || template.cookbook_name]
    cookbook.preferred_filename_on_disk_location(node, :templates, template.source)
  end
end
