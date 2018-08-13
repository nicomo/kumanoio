require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap-sass/assets/javascripts/bootstrap.js");
$(() => {

});


$("#flag-text").hover(
    function(){ $(this).addClass('btn-danger')},
    function(){ $(this).removeClass('btn-danger')}
)

$("#star-text").click(function(){
    $.ajax({
        type: 'POST',
        url: '/texts/' + $('a[id="star-text"]').attr('data-star-textid') + '/star',
        headers: {'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content')},
        success: function(){
            $("#glyph-star").removeClass('glyphicon-star-empty').addClass('glyphicon-star')
        }
    });
});    
/*  done(function() {
        console.log("req done");
    });

    fail(function(jqXHR, textStatus) {
        console.log("Request failed: " + textStatus);
    });
*/

