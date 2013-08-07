//Remove <ul> from Releases
$('.toctree-l1 > a').each(function(){
  if($(this).html() == "Releases") {
    $(this).parent().html('<a class="reference internal" href="/latest/releases/">Releases</a>');
  }
});

//Check URL
var path = window.location.href;
var pathName = path.split('/');
var num = pathName.length - 3;
var gettingstarted, terms, client, server, contributing;

switch(pathName[num]) {
  case 'gettingstarted':
    //case at Getting Started
    $('.toctree-l1 > a').each(function(){
      if($(this).html() == "Getting Started") {
        $(this).attr('state', 'open');
      }
    });
    break;
  case 'terms':
    //case at Terms
    $('.toctree-l1 > a').each(function(){
      if($(this).html() == "Terms") {
        $(this).attr('state', 'open');
      }
    });
    break;
  case 'client':
    //case at Client Reference
    $('.toctree-l1 > a').each(function(){
      if($(this).html() == "Client Reference") {
        $(this).attr('state', 'open');
      }
    });
    break;
  case 'server':
    //case at Server Reference
    $('.toctree-l1 > a').each(function(){
      if($(this).html() == "Server Reference") {
        $(this).attr('state', 'open');
      }
    });
    break;
  case 'contributing':
    //case at Contributing
    $('.toctree-l1 > a').each(function(){
      if($(this).html() == "Contributing") {
        $(this).attr('state', 'open');
      }
    });
    break;
  case 'releases':
    //code this out when releases gets filled out
    break;
  default:
}

