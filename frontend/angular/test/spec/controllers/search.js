/*global describe, beforeEach, expect, it, inject*/
'use strict';

describe('Controller: SearchCtrl', function () {

  // load the controller's module
  beforeEach(module('theDistributedBayApp'));

  var SearchCtrl,
    scope,
    httpBackend;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope, $httpBackend) {
    scope = $rootScope.$new();
    SearchCtrl = $controller('SearchCtrl', {
      $scope: scope
    });

    httpBackend = $httpBackend;
    httpBackend.when('GET', '/api/search?q=test').
      respond({ Torrents: [], Pages: 100 });
  }));

  it('should set initial variables correctly', function () {
    expect(scope.pages).toBe(0);
    expect(scope.page).toBe(0);
  });
  it('should should return currentSearchWithSort', function () {
    scope.sortDir = 'desc';
    scope.sort = 'Seeders';
    expect(scope.currentSearchWithSort('Age')).toBe('/?sort=Age:desc');

    expect(scope.currentSearchWithSort('Seeders')).toBe('/?sort=Seeders:asc');

    scope.sortDir = 'asc';
    expect(scope.currentSearchWithSort('Seeders')).toBe('/?sort=Seeders:desc');

    expect(scope.currentSearchWithSort('Age')).toBe('/?sort=Age:desc');
  });
  it('should set the page to 0 on query update and trigger a get request', function() {
    httpBackend.expectGET('/api/search?q=test');

    scope.page = 10;
    scope.query = 'test';
    scope.$apply();

    expect(scope.page).toBe(0);

    httpBackend.flush();
    expect(scope.pages).toBe(100);
  });
});
