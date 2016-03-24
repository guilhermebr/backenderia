angular.module('Backenderia')
.controller("PagesCreateController", function($http) {
    var controller = this;
    this.savePage = function(page) {
        controller.errors = null;
        $http({method: 'POST', url: '/api/page/', data: page})
        .catch(function(page){
            controller.errors = page.data.error;
        })
    };
})
