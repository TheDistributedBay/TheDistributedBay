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
    $scope.$watch('query', function(val) {
      var defer = data.search({
        q: val
      }).then(function(result) {
        console.log("Searh results", result);
        _.each(result, function(torrent) {
          torrent.Name = _.escape(torrent.Name);
        });
        $scope.torrents = result;
      }, function(err) {
      });
    });
  });
