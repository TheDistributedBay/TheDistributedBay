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
      controller: function($scope, $location, helpers) {
        $scope.headers = ['Age', 'Size', 'Seeders', 'Leechers'];
        $scope.sanitize = helpers.sanitizeName;
        var pieces = ($location.search().sort || '').split(':');
        $scope.sort = pieces[0];
        $scope.sortDir = pieces[1] || 'desc';
        $scope.sortBy = function(header) {
        };
        $scope.$watch('torrents', function() {
          _.each($scope.torrents, function(torrent) {
            torrent.TimeAgo = moment(torrent.CreatedAt).fromNow();
          });
        });
        $scope.humanFileSize = helpers.humanFileSize;
      }
    };
  });
