angular.module('bot-net',
    [
        'ngMaterial',
        'ngResource',
        'ui.router',
        'md.data.table'
    ])
    .config(function($stateProvider, $urlRouterProvider){
        $urlRouterProvider.otherwise('/');

        $stateProvider
            .state('talk', {
                url: '/',
                templateUrl: 'view/talk.html',
                controller: 'TalkController'
            })
            .state('detail', {
                url: '/:id',
                templateUrl: 'view/detail.html',
                controller: 'DetailController'
            })
            .state('start', {
                url: '/start/:talkName',
                templateUrl: 'view/start.html',
                controller: 'StartController'
            });
    });