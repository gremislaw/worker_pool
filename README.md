# worker_pool
## Тестовое задание в рамках vk internship.

- Worker-pool с возможностью динамически добавлять и удалять воркеры

## Cборка проекта

```bash
make build
```

## Запуск проекта

```bash
make run
```
В файле `/bin/out.txt` можно посмотреть действия программы

## Функционал

- Программа принимает на вход запросы, который принимает значения:
- - `add_job <num_of_jobs>` - создаст новые <strong>job</strong> в количестве <strong><_num_of_jobs_></strong>
- - `add_worker <id>` - cоздаст нового <strong>worker</strong> с <strong>id</strong> равным <strong><_id_></strong>
- - `delete_worker <id>` - удалит <strong>worker</strong> с <strong>id</strong> равным <strong><_id_></strong>
- - `stop` - принудительно завершит программу


## Пример

### Ввод

```bash
add_job 5
add_worker 3
add_worker 1
add_worker 2
delete_worker 3
add_job 4
```

### Вывод
```
worker 1 arrived.
worker 1 started job 1.
worker 3 arrived.
worker 3 started job 2.
worker 2 arrived.
worker 2 started job 3.
worker 2 finished job 3.
worker 2 started job 4.
worker 1 finished job 1.
worker 1 started job 5.
worker 3 finished job 2.
worker 3 kicked.
worker 1 finished job 5.
worker 1 started job 6.
worker 2 finished job 4.
worker 2 started job 7.
worker 1 finished job 6.
worker 1 started job 8.
worker 2 finished job 7.
worker 2 started job 9.
worker 2 finished job 9.
worker 1 finished job 8.
```