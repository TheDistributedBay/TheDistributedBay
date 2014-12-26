'use strict';

/**
 * @ngdoc function
 * @name theDistributedBayApp.controller:TorrentViewCtrl
 * @description
 * # TorrentViewCtrl
 * Controller of the theDistributedBayApp
 */
angular.module('theDistributedBayApp')
  .controller('TorrentViewCtrl', function ($scope, $routeParams, data, helpers) {
    $scope.sanitize = helpers.sanitizeName;
    $scope.humanFileSize = helpers.humanFileSize;
    data.getTorrent($routeParams.hash).then(function(torrent) {
      console.log("TORRENT", torrent);
      $scope.torrent = torrent;
    }, function(err) {
      alert('Not found:', err);
    });
  });
