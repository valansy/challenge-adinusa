package main

import (
	"log"
	"os"
	"os/signal"

	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	kingpin "github.com/alecthomas/kingpin/v2"

	"github.com/IBM/sarama"
)

var (
	brokerList        = kingpin.Flag("brokerList", "List of brokers to connect").Default("kafka:9092").Strings()
	topic             = kingpin.Flag("topic", "Topic name").Default("votes").String()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func main() {
	host := mustGetEnv("POSTGRES_HOST")
	port := mustGetEnv("POSTGRES_PORT")
	user := mustGetEnv("POSTGRES_USER")
	password := mustGetEnv("POSTGRES_PASSWORD")
	dbname := mustGetEnv("POSTGRES_DB")

	db := openDatabase(host, port, user, password, dbname)
	defer db.Close()

	pingDatabase(db)

	dropTableStmt := `DROP TABLE IF EXISTS votes`
	if _, err := db.Exec(dropTableStmt); err != nil {
		log.Panic(err)
	}

	createTableStmt := `CREATE TABLE IF NOT EXISTS votes (id VARCHAR(255) NOT NULL UNIQUE, vote VARCHAR(255) NOT NULL)`
	if _, err := db.Exec(createTableStmt); err != nil {
		log.Panic(err)
	}

	master := getKafkaMaster()
	defer master.Close()

	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				*messageCountStart++
				fmt.Printf("Received message: user %s vote %s\n", string(msg.Key), string(msg.Value))

				insertDynStmt := `INSERT INTO votes (id, vote) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET vote = $2`
				if _, err := db.Exec(insertDynStmt, string(msg.Key), string(msg.Value)); err != nil {
					log.Panic(err)
				}

			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed", *messageCountStart, "messages")
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func openDatabase(host, port, user, password, dbname string) *sql.DB {
    psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    for {
        db, err := sql.Open("postgres", psqlconn)
        if err == nil {
            return db
        }
        log.Printf("Failed to connect to database: %v. Retrying...", err)
        time.Sleep(2 * time.Second)
    }
}

func pingDatabase(db *sql.DB) {
	fmt.Println("Waiting for postgresql...")
	for {
		if err := db.Ping(); err == nil {
			fmt.Println("Postgresql connected!")
			return
		}
	}
}

func getKafkaMaster() sarama.Consumer {
	kingpin.Parse()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := *brokerList
	fmt.Println("Waiting for kafka...")
	for {
		master, err := sarama.NewConsumer(brokers, config)
		if err == nil {
			fmt.Println("Kafka connected!")
			return master
		}
	}
}
