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
            if ($location.path() !== '/search') {
              $location.path('/search');
            }
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
