angular.module('Backenderia')
.controller("PagesShowController", function($http, $routeParams) {
    var controller = this;
    $http({method: 'GET', url: '/api/page/' + $routeParams.id})
    .success(function(data){
        controller.page = data.data.page;
    });
});
