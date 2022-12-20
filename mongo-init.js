db = db.getSiblingDB("Wikime_test_performance")

db.createCollection('IdBase')
db.createCollection('Genres')
db.createCollection('Vk')
db.createCollection('Users')

db.Genres.insertOne({
    "_id": "Genres",
    "Genres": [
      "Гурман",
      "Космос",
      "Повседневность",
      "Сёдзё",
      "Исторический",
      "Фэнтези",
      "Сёдзё-ай",
      "Музыка",
      "Военное",
      "Этти",
      "Детектив",
      "Полиция",
      "Драма",
      "Ужасы",
      "Школа",
      "Сверхъестественное",
      "Работа",
      "Супер сила",
      "Безумие",
      "Экшен",
      "Комедия",
      "Вампиры",
      "Романтика",
      "Боевые искусства",
      "Юри",
      "Игры",
      "Пародия",
      "Эротика",
      "Спорт",
      "Триллер",
      "Сёнен-ай",
      "Смена пола",
      "Сёнен",
      "Самураи",
      "Детское",
      "Дзёсей",
      "Машины",
      "Хентай",
      "Сёдзе",
      "Приключения",
      "Сёдзе-ай",
      "Психологическое",
      "Гарем",
      "Яой",
      "Додзинси",
      "Меха",
      "Фантастика",
      "Сэйнэн",
      "Демоны",
      "Магия"
    ]
  })
db.IdBase.insertOne({
    "_id": "AnimeID",
    "LastId": NumberLong(0)
  })
db.IdBase.insertOne({
    "_id": "UserID",
    "LastId":  NumberLong(1)  
  })

db.Vk.insertOne({
  "_id": NumberLong(167298936),
  "InnerID":  NumberLong(0)
})

db.Users.insertOne({
  "_id":  NumberLong(0),
  "Nickname": "káйфøвàя дéвøчká❤",
  "Role": "root",
  "Favorites": [],
  "Watched": [],
  "Avatar": "/images/default_avatar.jpg",
  "Rated": []
})