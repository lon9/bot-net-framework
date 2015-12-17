angular.module('bot-net')
.factory('TweetService', function($resource){
    return $resource('/api/tweet/:id', null, {
        'update': {method: 'PUT'}
    });
});