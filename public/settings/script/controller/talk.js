angular.module('bot-net')
.controller('TalkController', function($scope, $mdDialog, $state, TalkService){
    var page = 1;
    var maxResults = 40;
    $scope.talks = TalkService.query({page:page, maxResults: maxResults});

    $scope.prev = function(){
       page--;
       $scope.talks = TalkService.query({page:page, maxResults:maxResults});
    };

    $scope.next = function(){
        page++;
        $scope.talks = TalkService.query({page:page, maxResults:maxResults});
    };

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