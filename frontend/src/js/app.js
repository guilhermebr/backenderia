angular.module('Backenderia', ['ngRoute'])
.config(['$routeProvider', function($routeProvider) {
    $routeProvider
        .when('/', {
            redirectTo: '/home'
        })
        .when('/home', {
            templateUrl: 'templates/home.html'
        })
        .when('/pages', {
            templateUrl: 'templates/apps/pages/index.html',
            controller: function ($http) {
                var controller = this;
                $http({method: 'GET', url: '/api/page/'}).success(function(data) {
                  controller.pages = data.data.pages;
                //   console.log(data.data)
                });
            },
            controllerAs: 'pageCtrl'
        })
        .when('/pages/new', {
            templateUrl: 'templates/apps/pages/new.html',
            controller: 'PagesCreateController',
            controllerAs: 'ctrl'
        })
        .when('/pages/:id', {
            templateUrl: 'templates/apps/pages/show.html',
            controller: 'PagesShowController',
            controllerAs: 'ctrl'
        })
        .otherwise({
            redirectTo: '/home'
        });
}]);
