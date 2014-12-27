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
    httpBackend.when('GET', '/api/search?q=test&sort=Seeders.Min:desc').
      respond({ Torrents: [], Pages: 100 });
    httpBackend.when('GET', '/api/search?q=test&p=1000&sort=CreatedAt:asc').
      respond({ Torrents: [], Pages: 100 });
  }));

  it('should set initial variables correctly', function () {
    expect(scope.pages).toBe(0);
    expect(scope.page).toBe(0);
  });

  it('should set the page to 0 on query update and trigger a get request', function() {
    httpBackend.expectGET('/api/search?q=test&sort=Seeders.Min:desc');

    scope.page = 10;
    scope.sort = 'Seeders';
    scope.sortDir = 'desc';
    scope.query = 'test';
    scope.$apply();

    expect(scope.page).toBe(0);

    httpBackend.flush();
    expect(scope.pages).toBe(100);
  });
  it('should update the scope variables on routeUpdate', inject(function($location) {
    httpBackend.expectGET('/api/search?q=test&p=1000&sort=CreatedAt:asc');
    $location.search({ p: 1000, q: 'test', sort: 'Age:asc' });
    scope.$emit('$routeUpdate');
    scope.$apply();
    httpBackend.flush();

    expect(scope.page).toBe(0);
    expect(scope.pages).toBe(100);
    expect(scope.query).toBe('test');
    expect(scope.sort).toBe('Age');
    expect(scope.sortDir).toBe('asc');
  }));
});
