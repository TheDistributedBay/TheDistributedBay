'use strict';

describe('Service: data', function () {

  // load the service's module
  beforeEach(module('theDistributedBayApp'));

  // instantiate service
  var data;
  beforeEach(inject(function (_data_) {
    data = _data_;
  }));
  it('should be able to execute searches with parameters', inject(function($httpBackend) {
    $httpBackend.when('GET', '/api/search?q=test&p=100&sort=Seeders:desc').
      respond({Torrents: ['duck'], Pages: 1020});
    $httpBackend.expectGET('/api/search?q=test&p=100&sort=Seeders:desc');
    var completed = false;
    data.search({
      q: 'test',
      p: '100',
      sort: 'Seeders:desc',
      moo: undefined
    }).then(function(resp) {
      completed = true;
      expect(resp.Pages).toEqual(1020);
      expect(resp.Torrents[0]).toEqual('duck');
    }, function(err) {
    });
    $httpBackend.flush();

    expect(completed).toEqual(true);
  }));

  it('should be able to get torrents by hash', inject(function($httpBackend) {
    $httpBackend.when('GET', '/api/torrent?hash=testhash').
      respond({Name: 'ducks'});
    $httpBackend.expectGET('/api/torrent?hash=testhash');
    var completed = false;
    data.getTorrent('testhash').then(function(resp) {
      completed = true;
      expect(resp.Name).toEqual('ducks');
    }, function(err) {
    });
    $httpBackend.flush();

    expect(completed).toEqual(true);
  }));

  it('should be able to add torrents', inject(function($httpBackend) {
    $httpBackend.when('POST', '/api/add_torrent').
      respond({Name: 'ducks'});
    var obj = {
      Name: 'ducks',
      MagnetLink: 'magnet:somelink',
      Description: 'woof',
      CategoryID: 10,
      CreatedAt: new Date(),
      Tags: ['banana']
    };
    $httpBackend.expectPOST('/api/add_torrent', obj);

    data.addTorrent(obj);
    $httpBackend.flush();
  }));
});
