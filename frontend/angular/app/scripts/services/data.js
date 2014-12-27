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
  .service('data', function ($http, $q, helpers) {
    var searchCancel;
    this.search = function(options) {
      if (searchCancel) {
        searchCancel.resolve();
      }
      options = _.merge({
        q: ''
      }, options);
      options = helpers.removeEmpty(options);
      var queryString = _.map(options, function(v, k) {
        return k + '=' + v;
      }).join('&');
      var defer = $q.defer();
      searchCancel = $q.defer();
      $http.get('/api/search?'+queryString, {timeout: searchCancel.promise, cache: true}).
        success(function(data, status, headers, config) {
          defer.resolve(data);
        }).
        error(function(data, status, headers, config) {
          defer.reject(data);
        });
      return defer.promise;
    };
    this.getTorrent = function(hash) {
      var defer = $q.defer();
      $http.get('/api/torrent?hash='+escape(hash), {cache: true}).
        success(function(data, status, headers, config) {
          defer.resolve(data);
        }).
        error(function(data, status, headers, config) {
          defer.reject(data);
        });
      return defer.promise;
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

      $http.post('/api/add_torrent', options).
        success(function(data, status, headers, config) {
          console.log('ADDED', data);
        }).
        error(function(data, status, headers, config) {
          alert('Error! ', data);
        });
    };
  });
