angular.module('bot-net')
.controller('TalkController', function($scope, $mdDialog, $state, $resource, $mdToast, TalkService){
    var page = 1;
    var maxResults = 40;
    $scope.talks = TalkService.query({page:page, maxResults: maxResults});

    //Prev page
    $scope.prev = function(){
       page--;
       $scope.talks = TalkService.query({page:page, maxResults:maxResults});
    };

    //Next page
    $scope.next = function(){
        page++;
        $scope.talks = TalkService.query({page:page, maxResults:maxResults});
    };

    // Add new talk
    $scope.newTalk = function(ev){
        $mdDialog.show({
            controller: NewTalkController,
            templateUrl: 'view/newtalk.html',
            parent: angular.element(document.body),
            parentEvent: ev,
            clickOutsideToClose:true
        })
        .then(function(newTalk){
           TalkService.save(newTalk, function(data){
                $state.go('detail', {
                    id: data.id
                });
           });
        });
    };

    // Delete tweets of the talk
    $scope.deleteTweet = function(talkId){
        var DelResource = $resource('/api/');
        DelResource.remove({
            talkId: talkId
        }, {}, function(data){
            $mdToast.show(
                $mdToast.simple()
                    .textContent('Deleted tweets')
                    .position({
                        top: true,
                        bottom: false,
                        right: true,
                        left: false
                    })
                    .hideDelay(3000)
            );
        });
    };
});

function NewTalkController($scope, $mdDialog){
    $scope.newTalk = {
        title:''
    };

    $scope.cancel = function(){
        $mdDialog.cancel();
    }

    $scope.post = function(talk){
        $mdDialog.hide(talk);
    }
}