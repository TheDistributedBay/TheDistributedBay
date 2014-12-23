'use strict';

/**
 * @ngdoc function
 * @name theDistributedBayApp.controller:SearchCtrl
 * @description
 * # SearchCtrl
 * Controller of the theDistributedBayApp
 */
angular.module('theDistributedBayApp')
  .controller('SearchCtrl', function ($scope, data) {
    $scope.torrents = [];
    $scope.pages = 0;
    $scope.$watch('query', function(val) {
        $scope.page = 0;
    });
    $scope.$watchGroup(['query', 'page'], function(val) {
      var defer = data.search({
        q: $scope.query,
        p: $scope.page
      }).then(function(result) {
        console.log("Search results", result);
        _.each(result.Torrents, function(torrent) {
          torrent.Name = _.escape(torrent.Name);
        });
        $scope.torrents = result.Torrents;
        $scope.pages = result.Pages;
        if ($scope.page >= $scope.pages) {
            $scope.page = 0;
        }
      }, function(err) {
      });
    });
  });
