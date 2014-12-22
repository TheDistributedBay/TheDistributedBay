'use strict';

describe('Controller: TorrentAddCtrl', function () {

  // load the controller's module
  beforeEach(module('theDistributedBayApp'));

  var TorrentAddCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    TorrentAddCtrl = $controller('TorrentAddCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
