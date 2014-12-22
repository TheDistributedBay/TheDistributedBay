'use strict';

/*global angular, _*/

/**
 * @ngdoc service
 * @name theDistributedBayApp.data
 * @description
 * # data
 * Service in the theDistributedBayApp.
 */
angular.module('theDistributedBayApp')
  .service('data', function ($http) {
    this.search = function(options) {
      options = _.merge({
      }, options);
    };
    this.addTorrent = function(options) {
      options = _.merge({
        MagnetLink: 'magnet:',
        Name: '',
        Description: '',
        CategoryID: 0,
        Size: 0,
        FileCount: 0,
        CreatedAt: new Date(),
        Tags: []
      }, options);

      console.log('ADD TORRENT', options);

      $http.post('/api/add_torrent', options).
        success(function(data, status, headers, config) {
          console.log('ADDED', data);
        }).
        error(function(data, status, headers, config) {
          alert('Error! ', data);
        });
    };
  });
