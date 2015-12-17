angular.module('bot-net')
.factory('TalkService', function($resource){
    return $resource('/api/talk/:id', null, {
        'update': {method: 'PUT'}
    });
});