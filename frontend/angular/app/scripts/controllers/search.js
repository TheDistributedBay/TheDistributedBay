'use strict';

/**
 * @ngdoc function
 * @name theDistributedBayApp.controller:SearchCtrl
 * @description
 * # SearchCtrl
 * Controller of the theDistributedBayApp
 */
angular.module('theDistributedBayApp')
  .controller('SearchCtrl', function ($scope, data, $location) {
    $scope.pages = 0;
    $scope.page = currentPage();
    $scope.$watch('query', function(val) {
      $scope.page = 0;
    });

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

    function currentPage() {
      return parseInt($location.search().p || 0, 10);
    }

    $scope.$on('$routeUpdate', function(){
      var search = $location.search();
      $scope.page = currentPage();
      $scope.query = search.q;
      window.scrollTo(window.scrollX,0)
    });
    $scope.$watchGroup(['query', 'page'], function() {
      $scope.loading = true;
      console.log('page', $scope.page);
      var defer = data.search({
        q: $scope.query,
        p: $location.search().p
      }).then(function(result) {
        console.log("Search results", result);
        _.each(result.Torrents, function(torrent) {
          torrent.Name = _.escape(torrent.Name);
        });
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

        $scope.pageLinks = _.times($scope.pages, function(page) {
          console.log(page, $scope.page, currentPage());
          return {
            link: currentSearchWithPage(page),
            page: page + 1,
            active: page === currentPage()
          };
        });
        $scope.loading = false;
      }, function(err) {
        console.log('SEARCH ERR', err);
      });
    });
  });
