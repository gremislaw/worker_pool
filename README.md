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
- - `add_job <data>` - создаст новый <strong>job</strong> с информацией <strong><_data_></strong>
- - `add_worker <id>` - cоздаст нового <strong>worker</strong> с <strong>id</strong> равным <strong><_id_></strong>
- - `delete_worker <id>` - удалит <strong>worker</strong> с <strong>id</strong> равным <strong><_id_></strong>
- - `stop` - принудительно завершит программу


## Пример

### Ввод

```bash
add_job backend
add_job frontend
add_job design
add_job product
add_job ml
add_worker 3
add_worker 3
add_worker 1
delete_worker 3
add_job backend
add_job frontend
add_job design
add_job product
# чутка подождем, пусть воркер 1 поработает один
add_worker 3
# подождите пока все выполнитcя
stop
```

### Вывод
```
Воркер 3 добавлен.
Воркер 3 обрабатывает строку: backend.
Воркер 1 добавлен.
Воркер 1 обрабатывает строку: frontend.
Воркер 3 обработал строку backend.
Воркер 3 обрабатывает строку: design.
Воркер 1 обработал строку frontend.
Воркер 1 обрабатывает строку: product.
Воркер 3 обработал строку design.
Воркер 3 удален.
Воркер 1 обработал строку product.
Воркер 1 обрабатывает строку: ml.
Воркер 1 обработал строку ml.
Воркер 1 обрабатывает строку: backend.
Воркер 1 обработал строку backend.
Воркер 1 обрабатывает строку: frontend.
Воркер 3 добавлен.
Воркер 3 обрабатывает строку: design.
Воркер 1 обработал строку frontend.
Воркер 1 обрабатывает строку: product.
Воркер 3 обработал строку design.
Воркер 1 обработал строку product.

```