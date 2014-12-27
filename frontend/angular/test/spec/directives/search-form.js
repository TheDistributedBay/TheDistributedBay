'use strict';

describe('Directive: searchForm', function () {

  // load the directive's module
  beforeEach(module('theDistributedBayApp'));

  var element,
    scope;

  beforeEach(module('templates'));

  beforeEach(inject(function ($rootScope) {
    scope = $rootScope.$new();
  }));

  it('should set initial variables correctly', inject(function ($compile, $location) {
    $location.search({
      q: 'test',
      p: 10,
      category: 'Anime,Software',
    });
    element = angular.element('<search-form></search-form>');
    element = $compile(element)(scope);
    scope.$digest();

    expect(scope.categories.length).toEqual(10);
    expect(scope.query).toEqual('test');
    expect(scope.page).toEqual(10);
    expect(scope.categories[1].checked).toEqual(true);
    expect(scope.categories[2].checked).toEqual(true);
  }));

  it('should zero the page count on a query change or category change', inject(function ($compile, $location) {
    $location.search({
      q: 'test',
      p: 10,
      category: 'Anime,Software',
    });
    element = angular.element('<search-form></search-form>');
    element = $compile(element)(scope);
    scope.$digest();

    expect(scope.query).toEqual('test');
    expect(scope.page).toEqual(10);

    scope.query = 'banana';
    scope.$apply();
    expect(scope.page).toEqual(0);

    scope.page = 10;
    scope.categories[3].checked = true;
    scope.checked(scope.categories[3]);
    scope.$apply();
    expect(scope.page).toEqual(0);
  }));

  it('should get bigger or smaller as the page demands', inject(function ($compile, $location) {
    element = angular.element('<div class="index-wrapper form-small"> <div class="index"> <search-form class="form-small"></search-form> </div> </div>');
    element = $compile(element)(scope);
    scope.$digest();

    $(element).appendTo(document.body);


    // Change to /
    scope.$emit('$routeChangeSuccess', {originalPath: '/'});
    expect($('search-form').hasClass('form-big')).toBeTruthy();
    expect($('.index-wrapper .index').hasClass('container')).toBeTruthy();

    scope.$emit('$routeChangeSuccess', {originalPath: '/search'}, {originalPath: '/'});
    expect($('search-form').hasClass('form-small')).toBeTruthy();
    expect($('.index-wrapper .index').hasClass('container')).toBeFalsy();

    $(element).remove();
  }));
});
