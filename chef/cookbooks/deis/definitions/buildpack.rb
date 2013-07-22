
define :buildpack, :name => nil, :git_url => nil, :target => nil do
  
  name = params[:name]
  git_url = params[:git_url]
  target = params[:target]
  
  if git_url == nil
    git_url = "git://github.com/heroku/heroku-buildpack-#{name}.git"
  end
  
  git target do
    repository git_url
    reference 'master'
    action :checkout
  end
  
end