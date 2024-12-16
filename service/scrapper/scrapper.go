package scrapper

import (
	"BotStocksScrapper/entity"
	sc "BotStocksScrapper/scrapper"
	"BotStocksScrapper/sender"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ScrapperService struct {
	stockScrapper sc.Scrapper
	stopChan      chan bool
	sender        *sender.Sender
	logger        entity.Logger
	db            any
	// TODO добавить поле сущность бд
}

func NewScrapperService(cfg entity.Config, tgClient *tgbotapi.BotAPI, chatID int64) (ScrapperService, error) {
	scrapper, err := sc.InitScrapper(cfg)
	if err != nil {
		return ScrapperService{}, err
	}

	sender := sender.NewSender(tgClient, chatID)

	return ScrapperService{
		stockScrapper: scrapper,
		sender:        sender,
		logger:        cfg.Logger,
		db:            nil,
	}, nil
}

// Блокирующая функция запускающая скраппер
func (s *ScrapperService) Scrap() error {
	stockChannel, err := s.stockScrapper.Scrape()
	if err != nil {
		s.logger.Errorf("Не удалось запустить скраппер: %s", err.Error())
		return err
	}

	for {
		select {
		case <-s.stopChan:
			s.stockScrapper.StopScrape()
			s.logger.Info("Остановлен сервис скраппера")
			return nil

		case stockInfo := <-stockChannel:
			// TODO Добавляем в бд сделку
			//  Получаем из бд все сделки к текущему моменту
			//  Дозаполняем StockInfo

			if stockInfo.IsAnomaly {
				err = s.sender.SendMsg(stockInfo)
				if err != nil {
					s.logger.Errorf("Ошибка отправки сообщения в канал: %s", err.Error())
				} else {
					s.logger.Infof("Успешно отправлена аномалия в чат: %v", stockInfo)
				}
			}
		}
	}
}

func (s *ScrapperService) Stop() {
	s.stopChan <- true
}
