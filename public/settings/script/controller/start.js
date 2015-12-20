angular.module('bot-net')
.controller('StartController', function($scope, $stateParams){
    $scope.talkName = $stateParams.talkName;
    $scope.tweets = [];
    var url = "ws://" + location.host + "/api/ws?talkName=" + $stateParams.talkName;
    ws = new WebSocket(url);
    ws.onmessage = function(event){
        var model = JSON.parse(event.data);
        $scope.tweets.push(model);
        $scope.$apply();
    };
});