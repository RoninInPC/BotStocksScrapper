package scrapper

import (
	"BotStocksScrapper/internal/hash"
	"BotStocksScrapper/internal/repository/logBase"
	"BotStocksScrapper/internal/repository/logBase/implLogBase"
	"BotStocksScrapper/internal/sender"
	"BotStocksScrapper/internal/sender/telegram"
	"time"

	"BotStocksScrapper/internal/entity"
	sc "BotStocksScrapper/internal/scrapper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ScrapperService struct {
	stockScrapper sc.Scrapper
	stopChan      chan bool
	sender        sender.Sender[entity.StockInfo]
	logger        entity.Logger
	baseLog       logBase.LogBase
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
		baseLog:       implLogBase.NewRedisRepository(implLogBase.NewRedisClient(cfg.RedisLog), hash.Nothing),
	}, nil
}

// Блокирующая функция запускающая скраппер
func (s *ScrapperService) Work() {
	stockChannel, err := s.stockScrapper.Scrape()
	if err != nil {
		s.logger.Errorf("Не удалось запустить скраппер: %s", err.Error())
		return
	}

	for {
		select {
		case <-s.stopChan:
			s.stockScrapper.StopScrape()
			s.logger.Info("Остановлен сервис скраппера")
			return

		case stockInfo, ok := <-stockChannel:
			if !ok {
				s.logger.Warn("Канал скраппера закрыт")
				return
			}
			if stockInfo.IsAnomaly {

				s.baseLog.Add(stockInfo.String())

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
