'use strict';

/**
 * @ngdoc function
 * @name theDistributedBayApp.controller:TorrentAddCtrl
 * @description
 * # TorrentAddCtrl
 * Controller of the theDistributedBayApp
 */
angular.module('theDistributedBayApp')
  .controller('TorrentAddCtrl', function ($scope, data) {
    $scope.addTorrent = function() {
      data.addTorrent({
        MagnetLink: $scope.magnet,
        Name: $scope.name,
        Description: $scope.description,
        Tags: _.map($scope.tags.split(','), function(tag) {
          return tag.trim();
        })
      });
    };
  });
