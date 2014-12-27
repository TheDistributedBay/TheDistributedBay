'use strict';

describe('Controller: TorrentViewCtrl', function () {

  // load the controller's module
  beforeEach(module('theDistributedBayApp'));

  var TorrentViewCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    TorrentViewCtrl = $controller('TorrentViewCtrl', {
      $scope: scope
    });
  }));

  /*it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });*/
});
