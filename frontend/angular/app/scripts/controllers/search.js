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
  .controller('SearchCtrl', function ($scope, data, $location) {

    function currentSearchWithPage(page) {
      var search = _.clone($location.search());
      if (page === 0) {
        delete search.p;
      } else {
        search.p = page;
      }
      return $location.path() + '?' + _.map(search, function(v, k) {
        return k + '=' + v;
      }).join('&');
    }
    $scope.currentSearchWithSort = function(sort) {
      var search = _.clone($location.search());
      search.sort = sort;
      delete search.p;

      return $location.path() + '?' + _.map(search, function(v, k) {
        return k + '=' + v;
      }).join('&');
    };

    function currentPage() {
      return parseInt($location.search().p || 0, 10);
    }

    $scope.pages = 0;
    $scope.page = currentPage();
    $scope.$watch('query', function(val) {
      $scope.page = 0;
    });

    $scope.$on('$routeUpdate', function(){
      var search = $location.search();
      $scope.page = currentPage();
      $scope.query = search.q;
      window.scrollTo(window.scrollX,0)
    });
    $scope.$watchGroup(['query', 'page', 'category'], function() {
      $scope.loading = true;
      var defer = data.search({
        q: $scope.query,
        p: $location.search().p,
        category: $location.search().category
      }).then(function(result) {
        console.log("Search results", result);
        $scope.torrents = result.Torrents;
        $scope.pages = result.Pages;
        if ($scope.page >= $scope.pages) {
            $scope.page = 0;
        }
        $scope.jumpLinks = _.mapValues({
          start: 0,
          prev: Math.max($scope.page - 1, 0),
          next: Math.min($scope.page + 1, $scope.pages - 1),
          end: $scope.pages - 1
        }, currentSearchWithPage);

        var curPage = currentPage();
        var start = Math.max(curPage - 5, 0);
        var end = Math.min(curPage + 5, $scope.pages - 1);
        $scope.pageLinks = [];
        var page;
        for (page = start; page <= end; page++) {
          $scope.pageLinks.push({
            link: currentSearchWithPage(page),
            page: page + 1,
            active: page === curPage
          });
        }
        $scope.loading = false;
      }, function(err) {
        console.log('SEARCH ERR', err);
      });
    });
  });
