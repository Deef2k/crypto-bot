package bot

import (
	"context"
	"sync"
	"testing"
)

func TestSubscriptionManager_AddSubscription(t *testing.T) {
	subscriber := SubscriptionManager{
		subscription: make(map[int64]context.CancelFunc),
		mutex:        sync.Mutex{},
	}
	_, cancel := context.WithCancel(context.Background()) //cancel-функция которая отменит контекст
	subscriber.mutex.Lock()
	subscriber.subscription[12345] = cancel //добавить подписку в мапу
	subscriber.mutex.Unlock()

	subscriber.mutex.Lock()
	_, exists := subscriber.subscription[12345] //пытаемся достать значение по ключу
	subscriber.mutex.Unlock()

	if !exists { //проверка на есть ли подписка ?
		t.Errorf("Подписка не была добавлена")
	}
}

func TestSubscriptionManager_DeleteSubscription(t *testing.T) {
	subscriber := SubscriptionManager{
		subscription: make(map[int64]context.CancelFunc),
		mutex:        sync.Mutex{},
	}
	_, cancel := context.WithCancel(context.Background()) //cancel-функция которая отменит контекст
	subscriber.mutex.Lock()
	subscriber.subscription[12345] = cancel //добавить подписку в мапу
	subscriber.mutex.Unlock()

	subscriber.mutex.Lock()
	delete(subscriber.subscription, 12345)
	subscriber.mutex.Unlock()

	subscriber.mutex.Lock()
	_, exists := subscriber.subscription[12345]
	subscriber.mutex.Unlock()

	if exists {
		t.Errorf("Подписка после удаления все еще существует")
	}
}

func TestSubscriptionManager_OverwriteSubscription(t *testing.T) {
	subscriber := SubscriptionManager{
		subscription: make(map[int64]context.CancelFunc),
		mutex:        sync.Mutex{},
	}
	ctx1, cancel1 := context.WithCancel(context.Background())
	_ = ctx1
	subscriber.mutex.Lock()
	subscriber.subscription[12345] = cancel1
	subscriber.mutex.Unlock()

	ctx2, cancel2 := context.WithCancel(context.Background())
	_ = ctx2

	subscriber.mutex.Lock()
	oldCancelFunc, exists := subscriber.subscription[12345]
	if exists {
		oldCancelFunc() //вызываем старую функцию отмены
	}
	subscriber.subscription[12345] = cancel2 //созроняем новую
	subscriber.mutex.Unlock()

	subscriber.mutex.Lock()
	currentCancel := subscriber.subscription[12345] //копируем подписку в переменную
	subscriber.mutex.Unlock()

	if currentCancel == nil {
		t.Error("Подписка не найдена в мапе")
	}

}
