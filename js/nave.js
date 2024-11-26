(function($) {
  "use strict"; // Start of use strict

  // Smooth scrolling using jQuery easing
  $('a[href*="#"]:not([href="#"])').click(function() {
      if (location.pathname.replace(/^\//, '') === this.pathname.replace(/^\//, '') && location.hostname === this.hostname) {
          var target = $(this.hash);
          target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
          if (target.length) {
              $('html, body').animate({
                  scrollTop: (target.offset().top)
              }, 1000, "easeInOutExpo");
              return false;
          }
      }
  });

  // Toggle menu on mobile view
  const menuBtn = document.querySelector('.menu-btn');
  const navbar = document.querySelector('.navbar');
  menuBtn.addEventListener('click', () => {
      navbar.classList.toggle('active');
  });

  // Close responsive menu when a scroll trigger link is clicked
  $('.navbar a').click(function() {
      $('.navbar').removeClass('active');
  });

  // Update active class on scroll
  $(window).on('scroll', function() {
      var scrollPos = $(document).scrollTop();
      $('.navbar a').each(function() {
          var currLink = $(this);
          var refElement = $(currLink.attr('href'));
          if (refElement.length) {
              if (refElement.position().top <= scrollPos && refElement.position().top + refElement.height() > scrollPos) {
                  $('.navbar a').removeClass('active');
                  currLink.addClass('active');
              } else {
                  currLink.removeClass('active');
              }
          }
      });
  });

  // Disable right-click
  document.addEventListener('contextmenu', function(e) {
      e.preventDefault();
  });

  // Disable keyboard shortcuts for inspect tools and copying
  document.addEventListener('keydown', function(e) {
      if (
          e.ctrlKey && (e.key === 'u' || e.key === 'c' || e.key === 'v' || e.key === 's' || e.key === 'p' || e.key === 'a') ||
          (e.ctrlKey && e.shiftKey && (e.key === 'i' || e.key === 'j')) ||
          (e.key === 'F12')
      ) {
          e.preventDefault();
      }
  });

  // Prevent text selection and copying
  document.addEventListener('selectstart', function(e) {
      e.preventDefault();
  });

  document.addEventListener('copy', function(e) {
      e.preventDefault();
  });

  // Disable dragging images or other elements
  document.addEventListener('dragstart', function(e) {
      e.preventDefault();
  });

})(jQuery); // End of use strict
