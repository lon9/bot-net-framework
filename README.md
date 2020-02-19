# Bot Net Framework

Bot Net Framework is not malware framework, it is "Twitter" Bot Net Framework.   
We can create network of Twitter bots. This makes easy to communicate with bots.   
All we have to do is just registering bots and configuring communications with Web UI, 
and then we can start communication on the web UI or CLI such as `curl`.

![example image1](https://bot.gyazo.com/2f08b8e61ecfff07f77f388dec11927a.gif "Exanple1")

## Usage

```
go get github.com/Rompei/bot-net-framework
go build
./bot-net-framework [OPTION]
```

Open browser and access the port number we decided in options.

And alse, we can start discussion on CLI   
Requesting URL `localhost:<port-number>/api/?talkName=<talk-name>`, and then bots will start talk.

## Optinons


```
-p, --port=      Port number
-k, --key=       Twitter consumer key
-s, --secret=    Twitter consumer secret
-d, --db=        Kind of database supported mysql, postgres, and sqlite
-o, --dboptions= Database options (See https://github.com/jinzhu/gorm#initialize-database)
```

And we can also configure these options from environment valuables.

```
BN_PORT                 Port number
BN_CONSUMER_KEY         Twitter consumer key
BN_CONSUMER_SECRET      Twitter consumer secret
BN_DATABASE             Kind of dat abase
BN_DB_OPTIONS           Database options
```
