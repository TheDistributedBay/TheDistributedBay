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
    this.humanFileSize = function(bytes, si) {
        var thresh = si ? 1000 : 1024;
        if(bytes < thresh) {
          return bytes + ' B';
        }
        var units = si ? ['kB','MB','GB','TB','PB','EB','ZB','YB'] : ['KiB','MiB','GiB','TiB','PiB','EiB','ZiB','YiB'];
        var u = -1;
        do {
            bytes /= thresh;
            ++u;
        } while(bytes >= thresh);
        return bytes.toFixed(1)+' '+units[u];
    };
  });
