angular.module('Backenderia')
.controller("NewsCreateController", function($http) {
    var controller = this;
    this.saveNews = function(news) {
        controller.errors = null;
        $http({method: 'POST', url: '/api/news/', data: news})
        .catch(function(news){
            controller.errors = news.data.error;
        })
    };
})
