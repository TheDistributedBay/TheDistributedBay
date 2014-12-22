'use strict';

/**
 * @ngdoc directive
 * @name theDistributedBayApp.directive:torrentDisplay
 * @description
 * # torrentDisplay
 */
angular.module('theDistributedBayApp')
  .directive('torrentDisplay', function () {
    return {
      templateUrl: 'views/torrent-display.html',
      restrict: 'E',
      link: function postLink(scope, element, attrs) {
        scope.torrents = [{
          name: '[Commie] Yowamushi Pedal Grande Road - 10 [C483199F].mkv',
          id: '13111114',
          magnet: 'magnet:?xt=urn:btih:46c860c2ed8a25140d7a6379d80f47a21459a64d&amp;dn=%5BCommie%5D+Yowamushi+Pedal+Grande+Road+-+10+%5BC483199F%5D.mkv&amp;xl=417382114&amp;dl=417382114&amp;tr=udp://open.demonii.com:1337/announce&amp;tr=udp://tracker.publicbt.com:80/announce&amp;tr=udp://tracker.istole.it:80/announce&amp;tr=http://open.nyaatorrents.info:6544/announce',
          date: new Date(),
          size: '398.05 MB',
          seeders: 106,
          leechers: 5,
          category: 'Anime'
        }];
        scope.sanitize = function(str) {
          return str.replace(/\W+/g, '-');
        };
      }
    };
  });
