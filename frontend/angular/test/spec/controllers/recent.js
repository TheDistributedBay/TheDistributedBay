'use strict';

describe('Controller: RecentCtrl', function () {

  // load the controller's module
  beforeEach(module('theDistributedBayApp'));

  var RecentCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    RecentCtrl = $controller('RecentCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
