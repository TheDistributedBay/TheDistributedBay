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
      scope: {
        torrents: '=',
        page: '=',
        pages: '=',
        noNav: '=',
        sort: '=',
        sortDir: '=',
        loading: '=',
        category: '='
      },
      controller: function($scope, $location, helpers) {
        $scope.headers = ['Age', 'Size', 'Seeders', 'Leechers'];
        $scope.sanitize = helpers.sanitizeName;
        $scope.$watch('torrents', function() {
          _.each($scope.torrents, function(torrent) {
            torrent.TimeAgo = moment(torrent.CreatedAt).fromNow();
          });
        });

        $scope.humanFileSize = helpers.humanFileSize;
        $scope.currentSearchWithSort = helpers.currentSearchWithSort;
        $scope.updateJumpLinks = function() {
          if ($scope.page >= $scope.pages) {
              $scope.page = 0;
          }
          var curPage = helpers.currentPage();
          $scope.jumpLinks = _.mapValues({
            start: 0,
            prev: Math.max(curPage - 1, 0),
            next: Math.min(curPage + 1, $scope.pages - 1),
            end: $scope.pages - 1
          }, helpers.currentSearchWithPage);
          var start = Math.max(curPage - 5, 0);
          var end = Math.min(curPage + 5, $scope.pages - 1);
          $scope.pageLinks = [];
          var page;
          for (page = start; page <= end; page++) {
            $scope.pageLinks.push({
              link: helpers.currentSearchWithPage(page),
              page: page + 1,
              active: page === curPage
            });
          }
        };
        if (!$scope.noNav) {
          $scope.updateJumpLinks();
          $scope.$watch('torrents', $scope.updateJumpLinks);
        }
      }
    };
  });
