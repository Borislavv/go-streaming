up:
	# запускает приложение через без сборки проекта
	docker-compose -f docker-compose.local.yaml up # подходит для локальной разработки

build:
	# запускает build прилоения, создает ротацию сборок, складывает их в ./cmd/builts
	docker-compose up # если при новой сборке возникла ошибка, будет запущена legacy или legacy_legacy
