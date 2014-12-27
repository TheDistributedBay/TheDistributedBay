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

  it('should create jump links correctly', inject(function ($compile) {
    element = angular.element('<torrent-display></torrent-display>');
    element = $compile(element)(scope);
    scope.$digest();
    scope.torrents = [{}];
    scope.pages = 10;
    scope.$apply();
    expect(scope.jumpLinks).toEqual({ start: '/?', prev: '/?', next: '/?p=1', end: '/?p=9' });
    var links = [ { link: '/?', page: 1, active: true}, { link: '/?p=1', page: 2, active: false}, { link: '/?p=2', page: 3, active: false}, { link: '/?p=3', page: 4, active: false}, { link: '/?p=4', page: 5, active: false}, { link: '/?p=5', page: 6, active: false} ]
    expect(_.pluck(scope.pageLinks, 'link')).toEqual(_.pluck(links, 'link'));
    expect(_.pluck(scope.pageLinks, 'page')).toEqual(_.pluck(links, 'page'));
    expect(_.pluck(scope.pageLinks, 'active')).toEqual(_.pluck(links, 'active'));
  }));
});
