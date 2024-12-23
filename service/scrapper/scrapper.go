package scrapper

import (
	"BotStocksScrapper/sender"
	"BotStocksScrapper/sender/telegram"
	"time"

	"BotStocksScrapper/entity"
	sc "BotStocksScrapper/scrapper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ScrapperService struct {
	stockScrapper sc.Scrapper
	stopChan      chan bool
	sender        sender.Sender[entity.StockInfo]
	logger        entity.Logger
	db            any
	// TODO добавить поле сущность бд
}

func NewScrapperService(cfg entity.Config, tgClient *tgbotapi.BotAPI, chatID int64) (ScrapperService, error) {
	scrapper, err := sc.InitScrapper(cfg)
	if err != nil {
		return ScrapperService{}, err
	}

	return ScrapperService{
		stockScrapper: scrapper,
		stopChan:      make(chan bool),
		sender:        telegram.NewSender(tgClient, chatID),
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

		case stockInfo, ok := <-stockChannel:
			if !ok {
				s.logger.Warn("Канал скраппера закрыт")
				return nil
			}
			if stockInfo.IsAnomaly {
				err = s.sender.Send(stockInfo)
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
	time.Sleep(2 * time.Second)
}
