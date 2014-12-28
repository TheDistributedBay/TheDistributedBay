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
    element = angular.element('<torrent-display></torrent-display>');
    element = $compile(element)(scope);
    scope.$digest();
    var childScope = element.children().scope();
    expect(childScope.headers).toEqual(['Age', 'Size', 'Seeders', 'Leechers']);
  }));

  it('should generate torrent relative times', inject(function ($compile) {
    scope.torrents = [{
      CreatedAt: new Date(),
    }];
    element = angular.element('<torrent-display torrents="torrents"></torrent-display>');
    element = $compile(element)(scope);
    scope.$digest();
    var childScope = element.children().scope();
    expect(childScope.torrents[0].TimeAgo).toEqual('a few seconds ago');
  }));

  it('should create jump links correctly', inject(function ($compile) {
    scope.pages = 10;
    element = angular.element('<torrent-display pages="pages" torrents="[{}]"></torrent-display>');
    element = $compile(element)(scope);
    scope.$digest();
    var childScope = element.children().scope();
    expect(childScope.jumpLinks).toEqual({ start: '/?', prev: '/?', next: '/?p=1', end: '/?p=9' });
    var links = [ { link: '/?', page: 1, active: true}, { link: '/?p=1', page: 2, active: false}, { link: '/?p=2', page: 3, active: false}, { link: '/?p=3', page: 4, active: false}, { link: '/?p=4', page: 5, active: false}, { link: '/?p=5', page: 6, active: false} ]
    expect(_.pluck(childScope.pageLinks, 'link')).toEqual(_.pluck(links, 'link'));
    expect(_.pluck(childScope.pageLinks, 'page')).toEqual(_.pluck(links, 'page'));
    expect(_.pluck(childScope.pageLinks, 'active')).toEqual(_.pluck(links, 'active'));
  }));
});
