//Remove <ul> from Releases
$('.toctree-l1 > a:contains("Releases")').siblings().remove();

//Check URL
var path = window.location.pathname.split('/');
var cleanedPath = $.grep(path,function(n){ return(n) });

// the second to last path segment is the section
// unless there's only 1 segment
if (cleanedPath.length == 1) {
    path = cleanedPath[0];
} else {
    path = cleanedPath[cleanedPath.length-2];
}

switch(path) {
  case 'understanding_deis':
    $('.toctree-l1 > a:contains("Understanding Deis")').attr('state', 'open');
    break;
  case 'installing_deis':
    $('.toctree-l1 > a:contains("Installing Deis")').attr('state', 'open');
    break;
  case 'using_deis':
    $('.toctree-l1 > a:contains("Using Deis")').attr('state', 'open');
    break;
  case 'managing_deis':
    $('.toctree-l1 > a:contains("Managing Deis")').attr('state', 'open');
    break;
  case 'troubleshooting_deis':
    $('.toctree-l1 > a:contains("Troubleshooting Deis")').attr('state', 'open');
    break;
  case 'customizing_deis':
    $('.toctree-l1 > a:contains("Customizing Deis")').attr('state', 'open');
    break;
  case 'roadmap':
    $('.toctree-l1 > a:contains("Roadmap")').attr('state', 'open');
    break;
  case 'contributing':
    $('.toctree-l1 > a:contains("Contributing")').attr('state', 'open');
    break;
  case 'reference':
  case 'client':
  case 'server':
  case 'terms':
    $('.toctree-l1 > a:contains("Reference Guide")').attr('state', 'open');
    break;
  default:
    $('.toctree-l1 > a:contains("Version")').attr('state', 'close');
    break;
}
