// sadly embedding this into broadcast2html.xslt doesn't work for Opera - 
// the 'if( now < dtstart )' ends up html escaped...

// moment.lang("de");
var canonical_url = ('' + window.location).replace(/\.xml$/,'');
$('.canonical-url').text( canonical_url );
$('.base-url').text( canonical_url.replace(/\/stations\/[^\/]+\/\d{4}\/\d{2}\/\d{2}\/(index|\d{4}%20.+)$/, '') );

var canonical_path = window.location.pathname.replace(/\.xml$/,'');
var dtstart = moment( $("meta[name='DC.format.timestart']").attr('content') );
var dtend = moment( $("meta[name='DC.format.timeend']").attr('content') );
var now = moment();
if( now < dtstart )
  $( 'html' ).addClass('is_future');
else if( now < dtend )
  $( 'html' ).addClass('is_current');
else
  $( 'html' ).addClass('is_past');

// display podcast links
var podasts_json_url = canonical_path + '.json';
$.ajax({ url: podasts_json_url, cache: true, dataType: 'json' }).success( function( data ) {
  // display mp3/enclosure dir link
  var enclosure_mp3_url = canonical_path.replace(/\/stations\//,'/enclosures/') + '.mp3';
  var enclosure_dir_url = enclosure_mp3_url.replace(/[^\/]+$/,'');
  $( 'a#enclosure_link' ).attr('href', enclosure_dir_url);
  $.ajax({ url: enclosure_mp3_url, type: 'HEAD', cache: true, }).success( function() {
    $( 'html' ).addClass('has_enclosure_mp3');
    $( 'a#enclosure_link' ).attr('href', enclosure_mp3_url);
    $( 'a#enclosure_link' ).attr('title', "Download: Rechte Maustaste + 'Speichern unter...'");
    $( '#enclosure audio source' ).attr('src', enclosure_mp3_url);
    $( '#enclosure' ).attr('style', 'display:block');
  });
  var has_ad_hoc = false;
  var names = data.podcasts.map( function(pc) {
    has_ad_hoc = has_ad_hoc || (pc.name == 'ad_hoc');
    return '<a href="../../../../../podcasts/' + pc.name + '/">' + pc.name + '</a>';
  } );
  $( '#podcasts' ).html( names.join(', ') );
  if( names.length == 0 ) {
    ;
  } else {
    $( 'p#enclosure' ).attr('style', 'display:block');
    $( 'html' ).addClass('has_podcast');
    if( has_ad_hoc ) {
      $( '#ad_hoc_action' ).attr('name', 'remove');
      $( '#ad_hoc_submit' ).attr('value', 'Nicht Aufnehmen');
    } else {
      $( '#ad_hoc_submit' ).attr('style', 'display:none');
    }
  }
});

// make date time display human readable
$( '#dtstart' ).html( moment(dtstart).format('ddd D[.] MMM YYYY HH:mm') );
$( '#dtend' ).html( moment(dtend).format('HH:mm') );

// add today/tomorrow links
$( '#prev_week' ).attr('href', '../../../' + moment(dtstart).subtract('days', 7).format() );
$( '#yesterday' ).attr('href', '../../../' + moment(dtstart).subtract('days', 1).format() );
$( '#tomorrow'  ).attr('href', '../../../' + moment(dtstart).add('days', 1).format() );
$( '#next_week' ).attr('href', '../../../' + moment(dtstart).add('days', 7).format() );

function finishAlldayCurrentEntry(a) {
  a.removeClass('is_past').addClass('is_current').append( jQuery('<span/>').text('jetzt') );
  // pastBC.append('<svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="150" height="150"><rect width="90" height="90" x="30" y="30" style="fill:#0000ff;fill-opacity:0.75;stroke:#000000"/></svg>');
}

// add all day broadcasts
$.ajax({ url: '.', type: 'GET', cache: true, dataType: 'xml', }).success( function(xmlBody) {
  var hasRecording = false;
  var pastBC = null;
  var allLinks = $(xmlBody).find('a').map( function() {
    var me = $(this);
    if( '../' == me.attr('href') )                    // ignore parent link
      return null;
    if( hasRecording )                                // previous entry was a .json recording marker
      me.addClass('has_podcast');
    if( hasRecording = me.attr('href').search(/\.json$/i) >= 0 ) // remember and swallow .json
      return null;
    var txt = me.text().replace(/\.xml$/, '');
    var ma = txt.match(/^(\d{2})(\d{2})\s+(.*?)$/);   // extract time and title
    if( ma ) {
      var t0 = dtstart.hours(ma[1]).minutes(ma[2]).seconds(0); // assumes same day
      me.attr('title', t0.format());
      me.text( t0.format('HH:mm') + ' ' + ma[3] );
      // set past/current/future class
      if( now < t0 ) {
        if(pastBC) {
          finishAlldayCurrentEntry(pastBC);
          pastBC = null;
        }
        me.addClass('is_future');
      } else {
        pastBC = me;
        me.addClass('is_past');
      }
    } else {
      me.text(txt);                                   // index usually.
    }
    me.attr('href', me.attr('href').replace(/\.xml$/, '') );  // make canonical url
    return this;
  });
  if( pastBC && now < dtstart.hours(24).minutes(0).seconds(0) )
    finishAlldayCurrentEntry(pastBC);
  $( '#allday' ).html( allLinks );
  $( '#allday a' ).wrap('<li>');
  $( '#allday' ).show();
});

// add whatsonnow station list
$.ajax({ url: '../../../../../app/now.lua', type: 'GET', cache: true, dataType: 'xml', }).success( function(xmlBody) {
  $( '#whatsonnow' ).html('');
  $(xmlBody).find( 'broadcast' ).map( function() {
    var bc = $(this);
    var title = bc.find("meta[name = 'DC.title']").attr('content');
    var iden  = bc.find("meta[name = 'DC.identifier']").attr('content');
    a = $('<a></a>').append( $('<span class="station">').text(iden.replace(/\/.*$/,'')) );
    a.attr('href', '../../../../' + iden);
    if( title )
      a.append('<br class="br"/>', $('<span class="broadcast"/>').text(title) );
    $( '#whatsonnow' ).append(a);
  });
  $( '#whatsonnow a' ).wrap( '<li>' );
});
