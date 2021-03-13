# Config service
- read from config file
- reload config file on change
- save config file at exit

# How to configure
Each service should expose config struct via GetConfig()
# config file format
{
    "reload": true,
    "save_on_exit": true,
    "services": {
        "database": {
            "database": "mydb",
            "username": "dbuser",
            "password": "mypass"
        },
        "kafka": {
            "brokers": ["somehost:9092", "anotherhost:9092"],
            "tls": true,
            "consumer": {
                "topics": ["topic1", "topic2"],
                "group": "mygroup"
                "workers":  3
            },
            "publisher": {
            }
        },
        "etl": {
            "entity": "pisces",
        },
        "maxmind": {
            "city": true,
            "isp": true
        },
    }
}

