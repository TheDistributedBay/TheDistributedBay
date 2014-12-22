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
        function updateSearch() {
          if ($scope.query) {
            $location.path('/search');
            $location.search('q=' + $scope.query);
          } else {
            $location.search('');
          }
        }
        $scope.$watch('query', function() {
          updateSearch();
        });
        $rootScope.$on('$routeChangeSuccess', function (event, current, previous) {
          if (current.originalPath === "/") {
            $('search-form.form-small').hide();
          } else {
            $('search-form.form-small').show();
          }
        });
      }
    };
  });
