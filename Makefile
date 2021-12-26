export MY_IP=0.0.0.0
up:
	docker-compose up -d
down:
	docker-compose down

TOPIC=foo
create_topic:
	docker run --net=host --rm confluentinc/cp-kafka:5.0.0 kafka-topics --create --topic ${TOPIC} --partitions 4 --replication-factor 2 --if-not-exists --zookeeper localhost:32181
