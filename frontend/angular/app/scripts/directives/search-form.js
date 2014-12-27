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
      controller: function ($scope, $rootScope, $location, helpers) {
        $scope.categories = [{
          name: 'All',
          checked: false,
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
        $scope.query = $location.search().q;
        $scope.page = $location.search().p || 0;
        var category = $location.search().category;
        if (category) {
          _.each(category.split(','), function(category) {
            _.each($scope.categories, function(cat) {
              if (category === cat.name) {
                cat.checked = true;
              }
            });
          });
        }
        function updateSearch() {
          var search = $location.search();
          search.category = _.reduce($scope.categories, function(result, category) {
            if (category.checked) {
              return result.concat(category.name);
            }
            return result;
          }, []).join(',');
          if (search.category === '' || search.category === 'All') {
            delete search.category;
          }
          $scope.category = search.category;
          if ($scope.query) {
            if ($location.path() !== '/search') {
              $location.path('/search');
            }
            search.q = $scope.query;

            // Strip page if set to zero
            if (parseInt($scope.page, 10) !== 0) {
              search.p = $scope.page;
            } else {
              delete search.p;
            }
          } else {
            delete search.q;
          }
          $location.search(search);
        }
        function monitorChange(cur, old) {
          if (cur !== old) {
            delete $location.search().p;
            $scope.page = 0;
          }
        }
        $scope.$watch('query', monitorChange);
        $scope.$watch('category', monitorChange);
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
        $scope.checked($scope.categories[1]);
        $scope.isRecent = helpers.isRecent;
        $scope.isSearch = function() {
          return !helpers.isRecent() && $location.path() === '/search';
        }
        $scope.isBrowse = function() {
          return $location.path() === '/browse';
        }
      }
    };
  });
