'use strict';

/*global angular, _*/

/**
 * @ngdoc function
 * @name theDistributedBayApp.controller:SearchCtrl
 * @description
 * # SearchCtrl
 * Controller of the theDistributedBayApp
 */
angular.module('theDistributedBayApp')
  .controller('SearchCtrl', function ($scope, data, $location, helpers) {

    $scope.isRecent = helpers.isRecent;

    $scope.pages = 0;
    $scope.page = helpers.currentPage();
    $scope.$watch('query', function(val) {
      $scope.page = 0;
    });

    var pieces = ($location.search().sort || '').split(':');
    $scope.sort = pieces[0];
    $scope.sortDir = pieces[1] || 'desc';

    $scope.$on('$routeUpdate', function(){
      var search = $location.search();
      $scope.query = search.q;
      $scope.page = helpers.currentPage();
      var pieces = ($location.search().sort || '').split(':');
      $scope.sort = pieces[0];
      $scope.sortDir = pieces[1] || 'desc';
      window.scrollTo(window.scrollX,0)
    });
    $scope.$watchGroup(['query', 'page', 'category', 'sort', 'sortDir'], function() {
      $scope.loading = true;
      var sort = '';
      console.log($scope.sort, $location.search());
      if ($scope.sort) {
        sort = ($scope.sort.replace(/Age/g, 'CreatedAt') + ':' + $scope.sortDir).
          replace(/Seeders:/g, 'Seeders.Min:').
          replace(/Leechers:/g, 'Leechers.Min:').
          replace(/Completed:/g, 'Completed.Min:');
      }
      data.search({
        q: $scope.query || '',
        p: $location.search().p || '',
        sort: sort,
        category: $location.search().category
      }).then(function(result) {
        console.log("Search results", result);
        $scope.torrents = result.Torrents;
        $scope.pages = result.Pages;

        $scope.loading = false;
      }, function(err) {
        console.log('SEARCH ERR', err);
      });
    });
  });
