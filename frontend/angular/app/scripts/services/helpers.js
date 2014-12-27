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

    /**
     * Removes empty elements from an object.
     *
     * @param {Object} the object to remove properties from.
     * @return {Object}
     */
    this.removeEmpty = function(object) {
      var obj = _.clone(object);
      _.each(obj, function(v, k) {
        if (v === undefined || v === '') {
          delete obj[k];
        }
      });
      return obj;
    };

    /**
     * Makes a name suitable for a URL by replacing all spaces and symbols with -.
     *
     * @param str {String} the string to sanitize
     * @return {String}
     */
    this.sanitizeName = function(str) {
      return (str || '').replace(/\W+/g, '-');
    };

    /**
     * Converts a number of bytes into a human readable number.
     *
     * @param bytes {Number} the number of bytes
     * @param si {Bool} an optional parameter to use base 1000 instead of 1024.
     * @return {String}
     */
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
