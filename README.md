# Casbin Kafka Watcher

[Casbin](https://github.com/casbin/casbin) watcher for kafka

> This watcher library will enable users to dynamically change
> casbin policies through kakfa messages

Influenced by https://github.com/rusenask/casbin-go-cloud-watcher

### Installation

```shell
go get -u github.com/wgarunap/casbin-kafka-watcher
```

### Usage

Configuration can be passed from the application without any strict 
env variable. To Initialize the watcher user need to parse
`kafka brokers` array , `topic name` and `callback function`.

On each restart all the messages from in the kafka topic will be synced 
to application policy and latest state will build on the application.

User can keep the kafka message `key`,`value` and `headers` with preferred 
encoding. such as `[]byte`, `string`, `json`,`avro`, `protobuf` or etc.

Message decoding logic and loading the policy into enforcer can
be done inside the callback function. The logic has to write manually 
by the user.

Example:
```go

```

### Suggested Kafka Topic Configuration
Compact the data for a same key and keep the latest message, enabling log compaction.
This config will keep the topic data without deleting.
```shell
log.cleanup.policy = compact
```

If you are sending whole casbin policy in a single kafka message, remember to set 
the `max.message.bytes` config in Kafka. It's default value is `1MB`

Keep the kafka topic partition count as 1. Otherwise, parallel updates through 
partitions will result in incorrect final policy state.

If you are running kafka as a cluster keep the replica count as 3 or more.

### Contact
Aruna Prabhashwara\
wg.aruna.p@gmail.com
