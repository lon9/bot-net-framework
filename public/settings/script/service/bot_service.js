angular.module('bot-net')
.factory('BotService', function($resource){
    return $resource('/api/bot/:id', null, {
        'update': {method: 'PUT'}
    });
});