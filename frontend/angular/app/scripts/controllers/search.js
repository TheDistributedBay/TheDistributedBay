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
      var dir = 'desc';
      if (sort === $scope.sort) {
        dir = $scope.sortDir === 'desc' ? 'asc' : 'desc';
      }
      var search = _.clone($location.search());
      search.sort = sort + ':' + dir;
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
      $scope.query = search.q;
      $scope.page = currentPage();
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
