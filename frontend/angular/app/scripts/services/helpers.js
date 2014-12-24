'use strict';

/**
 * @ngdoc service
 * @name theDistributedBayApp.helpers
 * @description
 * # helpers
 * Service in the theDistributedBayApp.
 */
angular.module('theDistributedBayApp')
  .service('helpers', function () {
    // AngularJS will instantiate a singleton by calling "new" on this function
    this.sanitizeName = function(str) {
      return (str || '').replace(/\W+/g, '-');
    };
  });
