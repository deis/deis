# Add some default users for sudo so testing doesn't end up in a
# broken system where the default user cannot sudo.
default['authorization']['sudo']['users'] = ['vagrant', 'ubuntu', 'ec2-user', 'root']
# make sure sudo is passwordless for tests.
default['authorization']['sudo']['passwordless'] = true
