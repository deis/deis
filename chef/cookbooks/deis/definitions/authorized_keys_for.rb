define :authorized_keys_for, :name => nil, :group => nil, :home => nil, :keys => [] do

  user  = params[:name]
  group = params[:group] || user
  home  = params[:home]  || "/home/#{user}"
  ssh_public_keys  = params[:keys]

  if ssh_public_keys.any?

    directory "#{home}/.ssh" do
      owner user
      group group
      mode 0700
      action :create
      only_if "test -d #{home}"
    end
    
    file "#{home}/.ssh/authorized_keys" do
      owner user
      group group
      content ssh_public_keys.join("\n")
    end
  end

end