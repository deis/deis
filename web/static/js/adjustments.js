function set_columns(is_server_reference) {
  is_server_reference = typeof is_server_reference !== 'undefined' ? is_server_reference : false;
  var margin = 0;
  var maxHeight = 0;

  //find the tallest column
  $('.column_calc').each(function() {
    if (maxHeight < $(this).height()) {
      maxHeight = $(this).height();
      console.log("M: " + maxHeight + " T: " + $(this).height());
    }
  });

  //511: height of the navigation. 96: height of the footer
  if (maxHeight > 923) {margin = maxHeight - 511 - 96;}
  if (is_server_reference == true){margin = margin + 80;}
  console.log("Max Height: " + maxHeight + " Margin: " + margin);
  
  //Set the margin above the footer
  $('.social-menu').css({'margin-top': (margin)});
}

//Close all accordions
$('.toctree-l1 > ul').hide();

//This variable checks if all accordions are closed. Used to ensure styling for Releases & FAQ page.
var allClosed;

//If a sub-item is currently being viewed, keep it's parent menu open
$('.toctree-l1 > a').each(function(){
  if($(this).attr('state') == 'open') {
    $(this).next().show();
    set_columns(); 
    allClosed = false;
    return false;
  } else {
    allClosed = true;
  }
});

if (allClosed == true) {set_columns();}


//If menu is closed when clicked, expand it
$('.toctree-l1 > a').click(function() {
  
  //Make the titles of open accordions dead links
  if ($(this).attr('state') == 'open') {return false;}

  //Clicking on a title of a closed accordion
  if($(this).attr('state') != 'open' && $(this).siblings().size() > 0) {
    $('.toctree-l1 > ul').hide();
    $('.toctree-l1 > a').attr('state', '');
    $(this).attr('state', 'open');
    $(this).next().slideDown();
    if($(this).html() == "Server Reference") {set_columns(true);} else {set_columns();}
    return false;
  }
});