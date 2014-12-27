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

  it('should be able to add torrents', inject(function($httpBackend) {
    $httpBackend.when('POST', '/api/add_torrent').
      respond({Name: 'ducks'});
    var obj = {
      Name: 'ducks',
      MagnetLink: 'magnet:somelink',
      Description: 'woof',
      CategoryID: 2,
      CreatedAt: new Date(),
      Tags: ['blah','foo']
    };
    $httpBackend.expectPOST('/api/add_torrent', function(resp) {
      var r = JSON.parse(resp);
      obj.CreatedAt = r.CreatedAt;
      expect(r).toEqual(obj);
      return true;
    });

    scope.name = obj.Name;
    scope.description = obj.Description;
    scope.magnet = obj.MagnetLink;
    scope.category = obj.CategoryID;
    scope.tags = 'blah,foo';
    scope.addTorrent();

    $httpBackend.flush();
  }));
});
