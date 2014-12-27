'use strict';

describe('Service: helpers', function () {

  // load the service's module
  beforeEach(module('theDistributedBayApp'));

  // instantiate service
  var helpers;
  beforeEach(inject(function (_helpers_) {
    helpers = _helpers_;
  }));

  it('should sanitize names correctly', function () {
    expect(helpers.sanitizeName('bla wf wububu 1992 -34043****D#$%^&*&^%$')).toBe('bla-wf-wububu-1992-34043-D-');
    expect(helpers.sanitizeName(null)).toBe('');
  });

  it('should humanize file lengths correctly', function () {
    expect(helpers.humanFileSize(0)).toBe('0 B');
    expect(helpers.humanFileSize(12345678)).toBe('11.8 MiB');
    expect(helpers.humanFileSize(12345678, true)).toBe('12.3 MB');
  });

  it('should remove empty elements from an object', function () {
    expect(helpers.removeEmpty({
      a: 'test',
      b: '',
      c: 1234,
      d: 0,
      e: undefined
    })).toEqual({
      a: 'test',
      c: 1234,
      d: 0
    });
  });

  it('should should return currentSearchWithSort', function () {
    expect(helpers.currentSearchWithSort('Age', 'Seeders', 'desc')).toBe('/?sort=Age:desc');
    expect(helpers.currentSearchWithSort('Seeders', 'Seeders', 'desc')).toBe('/?sort=Seeders:asc');
    expect(helpers.currentSearchWithSort('Seeders', 'Seeders', 'asc')).toBe('/?sort=Seeders:desc');
    expect(helpers.currentSearchWithSort('Age', 'Seeders', 'asc')).toBe('/?sort=Age:desc');
  });

  it('should correctly identify the recent route', inject(function($location) {
    expect(helpers.isRecent()).toBeFalsy();
    $location.path('/search');
    $location.search({ sort: "Age:desc" });
    expect(helpers.isRecent()).toBeTruthy();
  }));
});
