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
  .service('data', function ($http, $q) {
    this.search = function(options) {
      options = _.merge({
        q: ''
      }, options);
      var queryString = _.map(options, function(v, k) {
        return v + '=' + k;
      }).join('&');
      var defer = $q.defer();
      $http.get('/api/torrents?'+queryString).
        success(function(data, status, headers, config) {
          defer.resolve(data);
        }).
        error(function(data, status, headers, config) {
          defer.reject(data);
        });
      return defer;
    };
    this.addTorrent = function(options) {
      options = _.merge({
        MagnetLink: 'magnet:',
        Name: '',
        Description: '',
        CategoryID: 0,
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
