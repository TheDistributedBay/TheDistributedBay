'use strict';

/**
 * @ngdoc function
 * @name theDistributedBayApp.controller:BrowseCtrl
 * @description
 * # BrowseCtrl
 * Controller of the theDistributedBayApp
 */
angular.module('theDistributedBayApp')
  .controller('BrowseCtrl', function ($scope, data) {
    $scope.browseCategories = _.map(_.pluck($scope.categories, 'name').slice(1), function(category, i) {
      data.search({
        sort: "Seeders.Min:desc",
        size: 5,
        category: category
      }, true).then(function(result) {
        $scope.browseCategories[i].torrents = result.Torrents;
        $scope.browseCategories[i].loading = false;
      }, function(err) {
        console.log('SEARCH ERR', err);
      });
      return {
        name: category,
        torrents: [],
        loading: true,
      };
    });
  });
