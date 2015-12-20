angular.module('bot-net')
.controller('DetailController', function($scope, $stateParams, BotService, TalkService, TweetService){
    var talkId = $stateParams.id;
    $scope.talk = TalkService.get({id: talkId});

    // Getting tweets of the talk
    $scope.tweets = TweetService.query({talkId:talkId}, function(){
        $scope.tweets.push({
                id:0,
                text:"",
                talkId:parseInt(talkId, 10),
                sequence:$scope.tweets.length+1,
                bots:null
            });
    });

    // Getting bots
    $scope.bots = BotService.query();

    //Add tweet
    $scope.add = function(index, tweet){
        if(!validTweet(tweet)){return;}
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

    //Update tweet
    $scope.update = function(index, tweet){
        if(!validTweet(tweet)){return;}
        tweet.bot = null;
        tweet.botId = parseInt(tweet.botId);
        TweetService.update({}, tweet, function(data){
            $scope.tweets[index] = data;
        });
    };

    // Delete tweet
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

    //Validate tweet
    function validTweet(tweet){
        return !(tweet.botId == null || tweet.text == "");
    }


});
