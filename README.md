#Bot Net Framework

##Description
Bot Net Framework is not malware framework, it is "Twitter" Bot Net Framework.   
We can create network of Twitter bots. This makes easy to communicate with bots.   
All we have to do is just registering bots and configuring communications with Web UI, 
and then we can start communication on the web UI or CLI such as `curl`.

##Usage

```
go get github.com/Rompei/bot-net-framework
go build
./bot-net-framework [OPTION]
```

Open browser and access port number we decided in options.

##Optinons


```
  -p, --port=      Port number
  -k, --key=       Twitter consumer key
  -s, --secret=    Twitter consumer secret
  -d, --db=        Kind of database supported mysql, postgres, and sqlite
  -o, --dboptions= Database options (See https://github.com/jinzhu/gorm#initialize-database)
```

And we can configure these options from environment valuables.
