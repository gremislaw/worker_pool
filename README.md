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
- - `1 <count> <data>` либо `add_job <count> <data>` - создаст новые <strong>job'ы</strong> с информацией <strong><_data_></strong> в количестве <strong><_count_></strong>
- - `2 <count>` либо `add_worker <count>` - cоздаст новых <strong>worker'ов</strong> в количестве <strong>count</strong>
- - `3 <count>` либо `delete_worker <count>` - удалит <strong>worker'ов</strong> в количестве <strong>count</strong>
- - `4` либо `stop` - принудительно завершит программу

## Запуск unit-тестов

```bash
make test
```


## Пример

### Ввод

```bash
1 1000 backend
1 1000 frontend
1 100 design
1 100 product
1 10 ml
2 30
2 300
2 1000
3 300
1 100 backend
1 100 frontend
2 3
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
...
```