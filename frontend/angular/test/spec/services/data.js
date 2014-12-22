'use strict';

describe('Service: data', function () {

  // load the service's module
  beforeEach(module('theDistributedBayApp'));

  // instantiate service
  var data;
  beforeEach(inject(function (_data_) {
    data = _data_;
  }));

  it('should do something', function () {
    expect(!!data).toBe(true);
  });

});
