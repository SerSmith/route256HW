# Route256 2023

Домашний проект созданный в рамках курса route256 от ОЗОН.
https://route256.ozon.ru/go-developer

# Домашнее задание №1
- Создать скелеты трёх сервисов по описанию АПИ из файла contracts.md
- Структуру проекта сделать с учетом разбиения на слои, бизнес-логику писать отвязанной от реализаций клиентов и хендлеров
- Все хендлеры отвечают просто заглушками
- Сделать удобный враппер для сервера по тому принципу, по которому делали на воркшопе
- Придумать самостоятельно удобный враппер для клиента
- Все межсервисные вызовы выполняются. Если хендлер по описанию из contracts.md должен ходить в другой сервис, он должен у вас это успешно делать в коде.
- Общение сервисов по http-json-rpc
- должны успешно проходить make precommit и make run-all в корневой папке

# Домашнее задание №2
Перевести всё взаимодействие c сервисами на протокол gRPC.
Для этого:

- Использовать разделение на слои, созданное ранее, заменив слой HTTP на GRPC.
- Взаимодействие по HTTP полностью удалить и оставить только gRPC.
- В каждом проекте нужно добавить в Makefile команды для генерации кода из proto файла и установки нужных зависимостей.

Дополнительное задание на алмазик: добавить HTTP-gateway и proto-валидацию.

# Домашнее задание №3
Для каждого сервиса(где необходимо что-то сохранять/брать) поднять отдельную БД в docker-compose.
Сделать миграции в каждом сервисе (достаточно папки миграций и скрипта).
Создать необходимые таблицы.
Реализовать логику репозитория в каждом сервисе.
В качестве query builder-а можно использовать любую библиотеку (согласовать индивидуально с тьютором). Рекомендуется https://github.com/Masterminds/squirrel.
Драйвер для работы с postgresql: только pgx (pool).
В одном из сервисов сделать транзакционность запросов (как на воркшопе).

# Домашнее задание №4
Уменьшить время ответа Checkout.listCart при помощи worker pool. Запрашивать не более 5 SKU одновременно.
Worker pool нужно написать самостоятельно. Обязательное требование - читаемость кода и покрытие комментариями.
При общении с Product Service необходимо использовать лимит 10 RPS на клиентской стороне.
Допускается использование библиотечных рейт-лимитеров. В случае собственного читаемый код и комментарии обязательны.
Во всех слоях сервиса необходимо использовать контекст для возможности отмены вызова.

Задания на алмазик:

Синхронизировать рейт-лимитер при помощи БД.
Аннулировать заказы старше 10 минут в фоне. Позаботиться о том, чтобы реплики сервиса не штурмовали БД все вместе.


# Домашнее задание №5

Необходимо обеспечить полное покрытие бизнес-логики ручек ListCart или Purchase модульными тестами (go test -coverprofile).

Если вдруг ваши слои до сих пор не изолированны друг от друга через интерфейсы, необходимо это сделать.

В качестве генератора моков можете использовать, что душе угодно: mockery, minimoc, gomock, ... 

Задание на алмазик: добавить интеграционные тесты для проверки слоя взаимодействия с базой данных.

# Домашнее задание №6

LOMS пишет в Кафку изменения статусов заказов
Сервис нотификаций должен их вычитывать и отправлять нотификации об изменениях статуса заказа (писать в телегу)
Нотификация должна быть доставлена гарантированно и ровно один раз
Нотификации должны доставляться в правильном порядке

Алмазик: Весь новый функционал покрыт юнит тестами плюс сам код написан таким образом, что конфигурацию для нотификаций можно задавать в конфиге и прокидывать в основную логику

# Домашнее задание №7

1. Перевести логи всех сервисов на структурные. Советую zap, но можете использовать любой другой удобный пакет. Должна быть поддержка уровней логирования.

2. Покрыть все операции трейсами, настроить сбор трейсов в Jaeger, который должен подниматься в композе вместе со всеми сервисами. Должен быть доступ к веб-интерфейсу Jaeger-а.

3. Настроить отдачу сервисами метрик и их сбор Прометеем, запущенным в отдельном контейнере. Обязательно иметь счетчик запросов с детализацией по хендлерам и гистограмму времени обработки запросов с детализацией по кодам ответа и хендлерам для серверов, для клиентов только гистограмму. Сделать это лучше всего через интерсепторы. Плюсом можно добавить еще метрик от себя по вкусу, например для баз данных и запросам в апи продакт-сервиса.


# Домашнее задание №8
Задание
Нужно в сервис нотификаций добавить возможность получения истории уведомлений пользователя за конкретный период и закэшировать её:
1. Если пользователь уже запрашивал историю за конкретный период, то возвращать её из кэша
2. Для кэширования использовать либо LRU с воркшопа, либо любое другое решение, например, Redis
3. Функционал покрыть unit тестами
4. Требуется обеспечить простоту замены конкретной реализации кэша
