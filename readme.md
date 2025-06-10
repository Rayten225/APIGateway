# APIGateway

## Установка
Следуйте этим шагам, чтобы установить и запустить проект на вашем устройстве:

1. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/Rayten225/APIGateway
   ```
2. Перейдите в папку проекта:
   ```bash
   cd APIGateway
   ```
3. Соберите проект с помощью Docker:
   ```bash
   docker-compose up --build
   ```

## Использование
После установки вы можете проверить работоспособность этими командами 
   ```bash
   # 1. Создание новости через API Gateway
echo "1. Создание новости через API Gateway:"
curl -X POST http://localhost:8000/news -H "Content-Type: application/json" -d '{"title":"Тестовая новость","content":"Это тестовая новость"}'
echo -e "\n"

# 2. Получение списка новостей через API Gateway
echo "2. Получение списка новостей через API Gateway:"
curl http://localhost:8000/news
echo -e "\n"

# 3. Поиск новостей по названию через API Gateway
echo "3. Поиск новостей по названию через API Gateway:"
curl http://localhost:8000/news\?s\=Тестовая
echo -e "\n"

# 4. Получение деталей новости с ID=1 через API Gateway
echo "4. Получение деталей новости с ID=1 через API Gateway:"
curl http://localhost:8000/news/1
echo -e "\n"

# 5. Создание комментария через API Gateway
echo "5. Создание комментария через API Gateway:"
curl -X POST http://localhost:8000/comments -H "Content-Type: application/json" -d '{"news_id":1,"text":"Отличная новость!"}'
echo -e "\n"

# 6. Проверка текста на цензуру через API Gateway
echo "6. Проверка текста на цензуру через API Gateway:"
curl -X POST http://localhost:8000/censor -H "Content-Type: application/json" -d '{"text":"Отличная новость!"}'
echo -e "\n"

# 7. Создание комментария с запрещенным текстом
echo "7. Создание комментария с запрещенным текстом:"
curl -X POST http://localhost:8000/comments -H "Content-Type: application/json" -d '{"news_id":1,"text":"qwerty"}'
echo -e "\n"

# 8. Повторное получение деталей новости с ID=1 через API Gateway
echo "8. Повторное получение деталей новости:"
curl http://localhost:8000/news/1
echo -e "\n"

# 9. Получение списка новостей с request_id
echo "9. Получение списка новостей с request_id:"
curl http://localhost:8000/news\?request_id\=12345
echo -e "\n"

# 10. Проверка логов для request_id=12345     
echo "10. Проверка логов для request_id=12345:"
docker-compose logs api-gateway | grep 12345
echo -e "\n"
docker-compose logs news-service | grep 12345
echo -e "\n"                                     
docker-compose logs comment-service | grep 12345  
echo -e "\n"
docker-compose logs censor-service | grep 12345  
   ```

## Требования
- Docker, Docker-compose, Go
