'use strict';

/**
 * @ngdoc directive
 * @name theDistributedBayApp.directive:searchForm
 * @description
 * # searchForm
 */
angular.module('theDistributedBayApp')
  .directive('searchForm', function () {
    return {
      templateUrl: 'views/search-form.html',
      restrict: 'E',
      link: function postLink(scope, element, attrs) {
      },
      controller: function ($scope, $rootScope, $location) {
        $scope.query = $location.search().q;
        $scope.page = $location.search().p || 0;
        function updateSearch() {
          if ($scope.query) {
            if ($location.path() !== '/search') {
              $location.path('/search');
            }
            var search = $location.search();
            search.q = $scope.query;
            // Strip page if set to zero
            if (parseInt($scope.page, 10) !== 0) {
                search.p = $scope.page;
            } else {
                delete search.p;
            }
            $location.search(search);
          } else {
            $location.search('');
          }
        }
        $scope.$watchGroup(['query', 'page'], function() {
          updateSearch();
        });
        $rootScope.$on('$routeChangeSuccess', function (event, current, previous) {
          if (current.originalPath === "/") {
            $('.form-small').removeClass('form-small').addClass('form-big');
            $('.index-wrapper .index:not(.container)').addClass('container')
          } else if (previous && previous.originalPath === "/") {
            $('.form-big').removeClass('form-big').addClass('form-small');
            $('.index-wrapper .index.container').removeClass('container')
          }
        });
      }
    };
  });
