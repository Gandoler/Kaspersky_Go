## Тесты проекта

В этой папке собраны модульные и интеграционные тесты для компонентов проекта.

### Быстрый старт


1. Запустить тесты только для папки `TESTS`:

```bash
go test ./TESTS/...
```


### Что покрыто

- `TESTS/HandlersTESTS/handle_test.go`
  - Health-check: `GET /healthz` возвращает 200 OK
  - Enqueue: успешный `POST /enqueue` возвращает 202 Accepted и кладёт job в `state`
  - Enqueue (негатив): `GET /enqueue` → 405 Method Not Allowed
  - Enqueue (негатив): невалидный JSON → 400 Bad Request

- `TESTS/Adapters/memory_queue_test.go`
  - `Enqueue/Dequeue`: корректная запись/чтение
  - Переполнение: при полном буфере `Enqueue` возвращает ошибку
  - `Close`: после закрытия `Dequeue` возвращает `ok=false`

- `TESTS/Adapters/memory_state_store_test.go`
  - `Set/Get`: корректная запись/чтение состояния
  - Конкурентный доступ: параллельные `Set/Get` не падают (проверка thread-safety)

- `TESTS/Processors/retry_processor_test.go`
  - Успешная обработка: `Process` возвращает `done`
  - Ретраи: при `MaxRetries>0` возможен статус `queued` (проверяем вероятностно)

- `TESTS/WorkerPool/workerpool_test.go`
  - Интеграционный сценарий: `WorkerPool` забирает job из очереди и доводит до `done`

### Заметки о недетерминизме

`RetryProcessor` использует случайность и задержки (backoff). Для снижения флаки-тестов:

- В тесте на успех задаём сид `rand.Seed(1)`
- В тесте на ретраи делаем несколько попыток, чтобы поймать статус `queued`
- Таймауты в тестах подобраны с запасом, но остаются быстрыми
