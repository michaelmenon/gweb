package gweb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// WebStream ... manages the messages recieved for this stream
type webRedisStream struct {
	//redis client
	Rdb *redis.Client

	//the web ID of this service
	WebId string
}

// newRedisStream ... connect to a redis stream provided the redis host connection string
// conn string in the format `localhost:6379`
func newRedisStream(streamConn string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     streamConn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

// PostStreamMessage to default redis stream
func (w *webRedisStream) PostMessage(webId string, data string) error {
	b := map[string]interface{}{
		GWebRedisStreamEvent: data,
	}
	//create the redis stream if its not already there
	err := w.Rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream:     GWebRedisStream,
		NoMkStream: false,
		Values:     b,
	}).Err()

	return err

}

// ReadRedisStreamMessages ... read messages and return
// at a time 500 messages can be read
func (w *webRedisStream) ReadMessageStream() ([]GWebMessage, error) {

	// get the next message from the stream for the consumer group
	msgs, err := w.Rdb.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    GWebRedisStreamGroup,
		Consumer: w.WebId,
		Streams:  []string{GWebRedisStream, ">"},
		Count:    500,
	}).Result()
	streamMessages := make([]GWebMessage, 0)
	if err == nil {
		// process the messages
		for _, msg := range msgs {
			for _, entry := range msg.Messages {
				// handle the stream entry here
				// Send the message out to the channel "streamChannel" so the processing logic
				// can be handled by the module which calls the "ReadfromStream" function
				// For example, this can be called from "etabroadcast" module
				// streamChannel <- customEntry
				if streamData, ok := entry.Values[GWebRedisStreamEvent].(string); ok {
					customEntry := GWebMessage{Data: streamData, MessageId: entry.ID}
					streamMessages = append(streamMessages, customEntry)
				}

			}
		}
	}
	return streamMessages, err
}
