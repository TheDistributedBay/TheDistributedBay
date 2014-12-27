'use strict';

describe('Controller: TorrentViewCtrl', function () {

  // load the controller's module
  beforeEach(module('theDistributedBayApp'));

  var TorrentViewCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope, $location) {
    scope = $rootScope.$new();
    $location.path('/torrent/testhash/some-name');
    TorrentViewCtrl = $controller('TorrentViewCtrl', {
      $scope: scope,
      $routeParams: {
        hash: 'testhash',
      }
    });
  }));

  it('should make a search request and then display the information', inject(function($httpBackend) {
    var resp = {Name: 'ducks'};
    $httpBackend.when('GET', '/api/torrent?hash=testhash').
      respond(resp);
    $httpBackend.expectGET('/api/torrent?hash=testhash');
    $httpBackend.flush();

    expect(scope.torrent).toEqual(resp);
  }));
});
