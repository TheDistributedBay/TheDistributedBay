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
        scope.categories = [{
          name: 'All',
          checked: true,
        }, {
          name: 'Anime',
          checked: false,
        }, {
          name: 'Software',
          checked: false,
        }, {
          name: 'Games',
          checked: false,
        }, {
          name: 'Adult',
          checked: false,
        }, {
          name: 'Movies',
          checked: false,
        }, {
          name: 'Music',
          checked: false,
        }, {
          name: 'Other',
          checked: false,
        }, {
          name: 'Series & TV',
          checked: false,
        }, {
          name: 'Books',
          checked: false,
        }];
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
            search.category = _.reduce($scope.categories, function(result, category) {
              if (category.checked) {
                return result.concat(category.name);
              } else {
                return result;
              }
            }, []).join(',');
            if (search.category === '' || search.category === 'All') {
              delete search.category;
            }
            $scope.category = search.category;

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
        $scope.search = updateSearch;
        $scope.searchLucky = updateSearch;
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
        $scope.checked = function(box) {
          var numChecked = _.reduce($scope.categories, function(sum, category) {
            return sum + category.checked;
          }, 0);
          if (box.name === 'All') {
            _.each($scope.categories.slice(1), function(category) {
              category.checked = false;
            });
          } else {
            $scope.categories[0].checked = false;
          }
          if (numChecked === 0) {
            $scope.categories[0].checked = true;
          }
          updateSearch();
        }
      }
    };
  });
