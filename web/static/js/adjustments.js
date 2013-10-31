$(function(){
  var pageHeight;

  function set_columns() {
    pageHeight = 0;
    $('footer').css('margin-top', '0');
    $('.nav-border').css('height', '600');
    pageHeight = $(document).height();
    // console.log(pageHeight);
    $('.nav-border').css('height', pageHeight);

    var margin = pageHeight - 511 - 158;
    if (pageHeight < 800){
      margin = 270;
      $('body').css('height', '930');
      $('.nav-border').css('height', '930');
    }
    $('footer').css('margin-top', margin);
    $('footer').css('display', 'block');

    if($(window).width() > 1171){$('.docs-sidebar').css({'position': 'absolute', 'right': '0'})};
    if($(window).width() < 1171){$('.docs-sidebar').css('position', 'static');}
  }

//This variable checks if all accordions are closed. Used to ensure styling for Releases & FAQ page.
  var allClosed;

  //Close all accordions, besides the menu containing the page that you've clicked on.
  $('.toctree-l1').each(function(){
    if($(this).children('a').attr('state') == 'open') {
      $(this).children('ul').show();
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
      $(this).next().slideDown(function(){set_columns();});
      return false;
    }
  });

  $(window).resize(set_columns);
});
