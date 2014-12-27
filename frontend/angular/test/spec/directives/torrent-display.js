'use strict';

describe('Directive: torrentDisplay', function () {

  // load the directive's module
  beforeEach(module('theDistributedBayApp'));

  var element,
    scope;

  beforeEach(module('templates'));

  beforeEach(inject(function ($rootScope) {
    scope = $rootScope.$new();
  }));

  it('should configure the sorting correctly', inject(function ($compile, $location) {
    $location.search({
      sort: 'Age:asc',
    });
    element = angular.element('<torrent-display></torrent-display>');
    element = $compile(element)(scope);
    scope.$digest();
    expect(scope.headers).toEqual(['Age', 'Size', 'Seeders', 'Leechers']);
    expect(scope.sort).toEqual('Age');
    expect(scope.sortDir).toEqual('asc');
  }));

  it('should generate torrent relative times', inject(function ($compile) {
    element = angular.element('<torrent-display></torrent-display>');
    element = $compile(element)(scope);
    scope.$digest();
    scope.torrents = [{
      CreatedAt: new Date(),
    }];
    scope.$apply();
    expect(scope.torrents[0].TimeAgo).toEqual('a few seconds ago');
  }));
});
