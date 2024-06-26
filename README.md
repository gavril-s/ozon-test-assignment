# Тестовое задание для Ozon

Для ясности и упрощения проверки я оставлю небольшое описание того, что я тут написал.

### Типы запросов
Запросы (см. файл `graph/schema.graphqls`):
- `post` - получить пост по `id`.
- `postSnippets` - получить список постов со сниппетами контента (для списка постов не нужно загружать весь контент).
- `comments` - комментарии (отдельным запросом от поста). Параметры:
    - `postID` - `id` поста;
    - `parentID` - `id` родительского комментария;
    - `first` - количество комментариев, которые пользователь хочет получить (на первом уровне вложенности!);
    - `depth` - глубина вложенности (при `depth = 1` пользователь получит один слой комментариев без ответов);
    - `after` - `id` последнего комментария, который не интересует пользователя.
      
    **В общем, пользователь (как максимум) получает `first` комментариев, каждый из которых содержит `first` ответов, и т.д. - глубина всего этого дела `depth`.**
- `createPost` - создать пост.
- `createComment` - создать комментарий.

Общие параметры (на всякий случай):
- `first` - **везде** означает количество элементов, запрашиваемое клиентом.
- `after` - **везде** означает `id` последнего полученного клиентом элемента.

### Что я могу сказать в своё оправдание
Заранее отвечаю на вопросы, которые, как мне кажется, могут у Вас возникнуть.

1) Как Вы могли заметить, у меня 3 пакета с модельками

   Это может показаться немного запутывающим, и возможно есть немного более элегантное решение, но я считаю правильным разделить эти 3 вещи:
    - `graph/model` - сгенерированные автоматически модели того, что отдаётся в ответах пользователю;
    - `internal/storage/model` - модели данных, с которыми работает интерфейс `Storage`;
    - `internal/storage/db/model` - модели того, что хранится в базе данных.

    Конвертация между ними дешёвая, а плюсов такое разделение даёт много - можно полностью поменять структуру базы данных так, чтобы никто этого не заметил.

3) Возможно решение с фиксированным количеством показываемых ответов на каждый комментарий покажется Вам не самым лучшим.
 
    Я тоже так думаю, но, мне кажется, что оно, во-первых, не такое плохое, как кажется на первый взгляд (по ощущениям, что-то подобное делает Хабр, там ветки чуть ли не полностью раскрываются с самого начала), а во-вторых, легко исправимое. Единственное, что меня остановило от его исправления - у меня было очень мало времени на это задание из-за моей текущей занятости.
    А так, конечно, лучше было бы показывать меньше ответов с увеличением уровня вложенности.
