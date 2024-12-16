package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	intId, errAdd := store.Add(parcel)
	if intId != 0 || errAdd != nil {
		require.Error(t, err)
	}

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	parc, err := store.Get(intId)

	if err != nil {
		require.Error(t, err)
	}
	if parc != parcel {
		require.Error(t, err)
	}
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(intId)
	if err != nil {
		require.Error(t, err)
	}
	parc, err = store.Get(intId)
	if err == nil {
		require.Error(t, err)
	}
	defer db.Close()
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	intAdd, err := store.Add(parcel)
	if intAdd != 0 || err != nil {
		require.Error(t, err)
	}

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(intAdd, newAddress)
	if err != nil {
		require.Error(t, err)
	}
	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	parc, err := store.Get(intAdd)
	if parc.Address != newAddress {
		require.Error(t, err)
	}
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	intAdd, err := store.Add(parcel)
	if intAdd != 0 || err != nil {
		require.Error(t, err)
	}
	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := "I want live"
	err = store.SetStatus(intAdd, newStatus)
	if err != nil {
		require.Error(t, err)
	}
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	parc, err := store.Get(intAdd)
	if parc.Status != newStatus {
		require.Error(t, err)
	}
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcel)
		if id != 0 || err != nil {
			require.Error(t, err)
		} // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	if err != nil {
		require.Error(t, err)
	}
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	if len(storedParcels) != len(parcels) {

	}

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		if value, exit := parcelMap[parcel.Number]; exit {
			if value != parcel {
				require.Error(t, err)
			}
		} else {
			require.Error(t, err)
		}

	}
	require.NoError(t, err)
}
