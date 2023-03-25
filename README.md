# Apache Kafka & Go

## Инструменты для работы с Apache Kafka
1. https://github.com/redpanda-data/console
2. https://www.conduktor.io/
3. https://www.kafkatool.com/download.html
4. https://github.com/Shopify/sarama


## Consumer Group:
1. подов > партиций (часть подов не читает)
2. подов < партиций (часть подов читает больше партиций чем другие)
3. подов = партиций (каждый под читает свою партицию)