'use strict';

/**
 * @ngdoc overview
 * @name theDistributedBayApp
 * @description
 * # theDistributedBayApp
 *
 * Main module of the application.
 */
angular
  .module('theDistributedBayApp', [
    'ngAnimate',
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch'
  ])
  .config(function ($routeProvider, $locationProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl'
      })
      .when('/browse', {
        templateUrl: 'views/browse.html',
        controller: 'BrowseCtrl'
      })
      .when('/recent', {
        templateUrl: 'views/recent.html',
        controller: 'RecentCtrl'
      })
      .when('/torrent/add', {
        templateUrl: 'views/torrent/add.html',
        controller: 'TorrentAddCtrl'
      })
      .when('/search', {
        templateUrl: 'views/search.html',
        controller: 'SearchCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
    // use the HTML5 History API
    $locationProvider.html5Mode(true);
  });
