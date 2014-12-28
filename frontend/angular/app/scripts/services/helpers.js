'use strict';

/*global angular, _*/

/**
 * @ngdoc service
 * @name theDistributedBayApp.helpers
 * @description
 * # helpers
 * Service in the theDistributedBayApp.
 */
angular.module('theDistributedBayApp')
  .service('helpers', function ($location) {
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

    /**
     * Returns the current url but with a different page parameter.
     * @param {Number} the page number
     * @return {String}
     */
    this.currentSearchWithPage = function(page) {
      var search = _.clone($location.search());
      if (page === 0) {
        delete search.p;
      } else {
        search.p = page;
      }
      return $location.path() + '?' + _.map(search, function(v, k) {
        return k + '=' + v;
      }).join('&');
    };

    /**
     * Returns a url with the current path but with a different sort parameter.
     * @params {String} sort string
     * @return {String}
     */
    this.currentSearchWithSort = function(sort, curSort, sortDir, category) {
      var dir = 'desc';
      if (sort === curSort) {
        dir = sortDir === 'desc' ? 'asc' : 'desc';
      }
      var search = _.clone($location.search());
      search.sort = sort + ':' + dir;
      delete search.p;

      if (category) {
        search.category = category;
      }

      return '/search?' + _.map(search, function(v, k) {
        return k + '=' + v;
      }).join('&');
    };

    /**
     * Returns the current page from the location bar. Ex. "?p=0" -> 0
     * @return {Number}
     */
    this.currentPage = function() {
      return parseInt($location.search().p || 0, 10);
    };

    /**
     * Returns whether the current view is the recent one.
     * @return {Bool}
     */
    this.isRecent = function(){
      return $location.path() === '/search' && _.isEqual($location.search(), {sort: "Age:desc"});
    };
  });
