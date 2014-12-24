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
    data.getTorrent($routeParams.hash).then(function(torrent) {
      torrent.Name = _.escape(torrent.Name);
      torrent.Description = _.escape(torrent.Description);
      $scope.torrent = torrent;
    }, function(err) {
      alert('Not found:', err);
    });
  });
