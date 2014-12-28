'use strict';

describe('Controller: BrowseCtrl', function () {

  // load the controller's module
  beforeEach(module('theDistributedBayApp'));

  var BrowseCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    BrowseCtrl = $controller('BrowseCtrl', {
      $scope: scope
    });
  }));

  it('should configure categories and launch an http request', inject(function ($controller) {
    scope.categories = [{
      name: 'All',
      checked: false,
    }, {
      name: 'Anime',
      checked: false,
    }];
    scope.$apply();
    BrowseCtrl = $controller('BrowseCtrl', {
      $scope: scope
    });
    expect(scope.browseCategories[0].name).toEqual('Anime');
    expect(scope.browseCategories[0].torrents).toEqual([]);
    expect(scope.browseCategories[0].loading).toEqual(true);
  }));
});
