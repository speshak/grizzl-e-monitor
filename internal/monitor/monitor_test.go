package monitor

import (
	"testing"
	"time"

	"github.com/speshak/grizzl-e-prom/pkg/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConnectAPI struct {
	mock.Mock
}

func (m *MockConnectAPI) GetStations() ([]connect.Station, error) {
	args := m.Called()
	return args.Get(0).([]connect.Station), args.Error(1)
}

func (m *MockConnectAPI) AssertValidToken() error {
	// Just short cut this
	return nil
}

func (m *MockConnectAPI) Login() error {
	// Just short cut this too
	return nil
}

func (m *MockConnectAPI) Logout() error {
	// Just short cut this too
	return nil
}

func (m *MockConnectAPI) SetDebug() {}

func (m *MockConnectAPI) GetTransactionStatistics(stationID string) (connect.TransactionStats, error) {
	args := m.Called(stationID)
	return args.Get(0).(connect.TransactionStats), args.Error(1)
}

func (m *MockConnectAPI) GetStation(stationID string) (connect.Station, error) {
	args := m.Called(stationID)
	return args.Get(0).(connect.Station), args.Error(1)
}

func (m *MockConnectAPI) GetAllTransactions(stationID string) ([]connect.Transaction, error) {
	args := m.Called(stationID)
	return args.Get(0).([]connect.Transaction), args.Error(1)
}

func (m *MockConnectAPI) GetTransaction(transactionID string) (connect.Transaction, error) {
	args := m.Called(transactionID)
	return args.Get(0).(connect.Transaction), args.Error(1)
}

func (m *MockConnectAPI) GetTransactions(stationId string, limit int, offset int) ([]connect.Transaction, error) {
	args := m.Called(stationId, limit, offset)

	//return args.Get(0).(connect.Transaction), args.Error(1)
	return args.Get(0).([]connect.Transaction), args.Error(1)
}

type MockTransactionHistoryPublisher struct {
	mock.Mock
}

func (m *MockTransactionHistoryPublisher) TransactionPublished(transactionID string) bool {
	args := m.Called(transactionID)
	return args.Bool(0)
}

func (m *MockTransactionHistoryPublisher) PublishTransactionHistory(stationID string, transaction connect.Transaction) {
	m.Called(stationID, transaction)
}

type MockTransactionStatsPublisher struct {
	mock.Mock
}

func (m *MockTransactionStatsPublisher) PublishTransactionStats(stationID string, stats connect.TransactionStats) {
	m.Called(stationID, stats)
}

type MockStationStatusPublisher struct {
	mock.Mock
}

func (m *MockStationStatusPublisher) PublishStationStatus(station connect.Station) {
	m.Called(station)
}

func TestMonitorStations(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetStations").Return([]connect.Station{{ID: "station1"}}, nil)

	monitor := &StationMonitor{
		Connect: mockConnectAPI,
	}

	go func() {
		time.Sleep(2 * time.Second)
	}()

	err := monitor.MonitorStations()
	assert.NoError(t, err)
	mockConnectAPI.AssertExpectations(t)
}

func TestMonitorStation(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetTransactionStats", "station1").Return(connect.TransactionStats{}, nil)
	mockConnectAPI.On("GetStation", "station1").Return(connect.Station{ID: "station1"}, nil)
	mockConnectAPI.On("GetAllTransactions", "station1").Return([]connect.Transaction{{ID: "trans1"}}, nil)
	mockConnectAPI.On("GetTransaction", "trans1").Return(connect.Transaction{ID: "trans1"}, nil)

	mockTransactionHistoryPublisher := new(MockTransactionHistoryPublisher)
	mockTransactionHistoryPublisher.On("TransactionPublished", "trans1").Return(false)
	mockTransactionHistoryPublisher.On("PublishTransactionHistory", "station1", mock.Anything)

	mockTransactionStatsPublisher := new(MockTransactionStatsPublisher)
	mockTransactionStatsPublisher.On("PublishTransactionStats", "station1", mock.Anything)

	mockStationStatusPublisher := new(MockStationStatusPublisher)
	mockStationStatusPublisher.On("PublishStationStatus", mock.Anything)

	monitor := &StationMonitor{
		Connect:                     mockConnectAPI,
		TransactionHistoryPublisher: mockTransactionHistoryPublisher,
		TransactionStatsPublisher:   mockTransactionStatsPublisher,
		StationStatusPublisher:      mockStationStatusPublisher,
	}

	station := connect.Station{ID: "station1"}
	monitor.MonitorStation(station)

	mockConnectAPI.AssertExpectations(t)
	mockTransactionHistoryPublisher.AssertExpectations(t)
	mockTransactionStatsPublisher.AssertExpectations(t)
	mockStationStatusPublisher.AssertExpectations(t)
}

func TestTransactionStats(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetTransactionStats", "station1").Return(connect.TransactionStats{}, nil)

	mockTransactionStatsPublisher := new(MockTransactionStatsPublisher)
	mockTransactionStatsPublisher.On("PublishTransactionStats", "station1", mock.Anything)

	monitor := &StationMonitor{
		Connect:                   mockConnectAPI,
		TransactionStatsPublisher: mockTransactionStatsPublisher,
	}

	station := connect.Station{ID: "station1"}
	monitor.transactionStats(station)

	mockConnectAPI.AssertExpectations(t)
	mockTransactionStatsPublisher.AssertExpectations(t)
}

func TestStationStats(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetStation", "station1").Return(connect.Station{ID: "station1"}, nil)

	mockStationStatusPublisher := new(MockStationStatusPublisher)
	mockStationStatusPublisher.On("PublishStationStatus", mock.Anything)

	monitor := &StationMonitor{
		Connect:                mockConnectAPI,
		StationStatusPublisher: mockStationStatusPublisher,
	}

	station := connect.Station{ID: "station1"}
	monitor.stationStats(station)

	mockConnectAPI.AssertExpectations(t)
	mockStationStatusPublisher.AssertExpectations(t)
}

func TestTransactionHistory(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetAllTransactions", "station1").Return([]connect.Transaction{{ID: "trans1"}}, nil)
	mockConnectAPI.On("GetTransaction", "trans1").Return(connect.Transaction{ID: "trans1"}, nil)

	mockTransactionHistoryPublisher := new(MockTransactionHistoryPublisher)
	mockTransactionHistoryPublisher.On("TransactionPublished", "trans1").Return(false)
	mockTransactionHistoryPublisher.On("PublishTransactionHistory", "station1", mock.Anything)

	monitor := &StationMonitor{
		Connect:                     mockConnectAPI,
		TransactionHistoryPublisher: mockTransactionHistoryPublisher,
	}

	station := connect.Station{ID: "station1"}
	monitor.transactionHistory(station)

	mockConnectAPI.AssertExpectations(t)
	mockTransactionHistoryPublisher.AssertExpectations(t)
}
