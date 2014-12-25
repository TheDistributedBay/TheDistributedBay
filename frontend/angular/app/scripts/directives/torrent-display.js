'use strict';

/**
 * @ngdoc directive
 * @name theDistributedBayApp.directive:torrentDisplay
 * @description
 * # torrentDisplay
 */
angular.module('theDistributedBayApp')
  .directive('torrentDisplay', function (helpers) {
    return {
      templateUrl: 'views/torrent-display.html',
      restrict: 'E',
      controller: function($scope, $location) {
        $scope.headers = ['Age', 'Size', 'Seeders', 'Leechers'];
        $scope.sanitize = helpers.sanitizeName;
        $scope.sort = $location.search().sort;
        $scope.sortBy = function(header) {
          $scope.sort = header;
        };
        $scope.$watch('torrents', function() {
          _.each($scope.torrents, function(torrent) {
            torrent.TimeAgo = moment(torrent.CreatedAt).fromNow();
          });
        });
      }
    };
  });
