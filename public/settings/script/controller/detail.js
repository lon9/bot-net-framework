angular.module('bot-net')
.controller('DetailController', function($scope, $stateParams, BotService, TalkService, TweetService){
    var talkId = $stateParams.id;
    $scope.talk = TalkService.get({id: talkId});
    $scope.tweets = TweetService.query({talkId:talkId}, function(){
        $scope.tweets.push({
                id:0,
                text:"",
                talkId:parseInt(talkId, 10),
                sequence:$scope.tweets.length+1,
                bots:null
            });
    });
    $scope.bots = BotService.query();

    $scope.add = function(index, tweet){
        tweet.botId = parseInt(tweet.botId);
        TweetService.save({}, tweet, function(data){
            $scope.tweets[index] = data;
            $scope.tweets.push({
                id:0,
                text:"",
                talkId:parseInt(talkId, 10),
                sequence:$scope.tweets.length+1,
                bots:null
            });
        });
    };

    $scope.update = function(index, tweet){
        TweetService.update({}, tweet, function(data){
            $scope.tweets[index] = data;
        });
    };

    $scope.delete = function(idnex, tweet){
        TweetService.remove({id:tweet.id}, {}, function(){
            $scope.tweets = TweetService.query({talkId:talkId}, function(){
                $scope.tweets.push({
                                id:0,
                                text:"",
                                talkId:$scope.talk.id,
                                sequence:$scope.tweets.length,
                                bots:null
                            });
            });
        });
    };


});
