# Описание
#### Telegram Бот-помошник, который за тебя пытается найти перевод и сделать анализ слов и строк получая информацию парся сторонние ресурсы
#### Такие ресурсы как 
**[dictionary.cambridge.org](https://dictionary.cambridge.org)** 

**[multitran.com](https://multitran.com)** 

**[wooordhunt.ru](https://wooordhunt.ru)** 

#### Умеет отдавать информацию разными способами, картинки, аудиосообщения
#### Также встроена база(mysql) куда пишутся всё что мы искали и делается кеш(redis) сырых/полусырых данных, того что уже искали

# Демо
![Демо](https://github.com/efremovigor/telegram-bot-golang/blob/master/demo.gif)

# Установка
```
cp .env_example .env
docker-compose up -d
```
