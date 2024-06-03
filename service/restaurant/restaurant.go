package service

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"restaurant-service/models"

	"restaurant-service/service"
)

type restaurantSvc struct {
	restaurantStore service.Restaurant
	kProducer       sarama.SyncProducer
}

func New(restaurantStore service.Restaurant, kProducer sarama.SyncProducer) restaurantSvc {
	return restaurantSvc{restaurantStore: restaurantStore, kProducer: kProducer}
}

func (s *restaurantSvc) Create(payload *models.Restaurant) error {
	err := s.restaurantStore.Create(payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *restaurantSvc) Read() ([]*models.Restaurant, error) {
	return s.restaurantStore.Read()
}

func (s *restaurantSvc) Update(restaurant *models.Restaurant) error {
	return s.restaurantStore.Update(restaurant)
}

func (s *restaurantSvc) OrderStatusUpdate(restaurant *models.OrderStatusUpdate) error {
	if err := s.restaurantStore.UpdateOrderStatus(restaurant.OrderID, restaurant.Status); err != nil {
		return err
	}

	rest, err := s.restaurantStore.GetByID(restaurant.RestaurantID)
	if err != nil {
		return err
	}

	return s.publishOrderStatusUpdatedEvent(restaurant.OrderID, restaurant.Status, rest.Latitude, rest.Longitude)

	//return r.restaurantStore.Update(restaurant)
}

func (s *restaurantSvc) Delete(id int) error {
	return s.restaurantStore.Delete(id)
}

func (s *restaurantSvc) ProcessOrderPlacedEvent(message []byte) error {
	var event models.OrderPlacedEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return err
	}

	if err := s.restaurantStore.InsertOrder(&event); err != nil {
		return err
	}

	// Publish order status updated event to Kafka
	return nil
}

func (s *restaurantSvc) publishOrderStatusUpdatedEvent(orderID int, status string, lat, long float64) error {
	event := models.OrderStatusUpdatedEvent{
		OrderID: orderID,
		Status:  status,
		Lat:     lat,
		Long:    long,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	kafkaMessage := &sarama.ProducerMessage{
		Topic: "order-status-updated",
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := s.kProducer.SendMessage(kafkaMessage)
	if err != nil {
		return err
	}

	fmt.Printf("Order status update sent to partition %d at offset %d\n", partition, offset)
	return nil
}
