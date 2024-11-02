package monitor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/speshak/grizzl-e-monitor/pkg/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	gocronmocks "github.com/go-co-op/gocron/mocks/v2"
)

// Create mocks for the ConnectAPI, TransactionHistoryPublisher, TransactionStatsPublisher, and StationStatusPublisher
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

func (m *MockConnectAPI) ParseToken() (*jwt.Token, connect.TokenClaims, error) {
	// Just short cut this too
	return nil, connect.TokenClaims{}, nil
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

	return args.Get(0).([]connect.Transaction), args.Error(1)
}

type MockTransactionHistoryPublisher struct {
	mock.Mock
}

func (m *MockTransactionHistoryPublisher) Close() error {
	return nil
}

func (m *MockTransactionHistoryPublisher) TransactionPublished(transaction connect.Transaction) bool {
	args := m.Called(transaction)
	return args.Bool(0)
}

func (m *MockTransactionHistoryPublisher) PublishTransactionHistory(stationID string, transaction connect.Transaction) error {
	m.Called(stationID, transaction)
	return nil
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

func (m *MockStationStatusPublisher) Close() error {
	return nil
}

type MockSingleStationMonitor struct {
	mock.Mock
}

func (m *MockSingleStationMonitor) MonitorStation(ctx context.Context, station connect.Station) {
	m.Called(ctx, station)
}

func TestMonitorConstructor(t *testing.T) {
	monitor := NewStationMonitor(&Config{
		APIHost:  "https://example.com",
		Username: "myUser",
		Password: "myPass",
		Debug:    true,
	})

	assert.NotNil(t, monitor)
	assert.NotNil(t, monitor.Connect)
}

func TestTransactionStats(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetTransactionStatistics", "station1").Return(connect.TransactionStats{}, nil)

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
	mockTransactionHistoryPublisher.On("TransactionPublished", connect.Transaction{ID: "trans1"}).Return(false)
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

func TestExistingTransactionHistory(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetAllTransactions", "station1").Return([]connect.Transaction{{ID: "trans1"}}, nil)

	mockTransactionHistoryPublisher := new(MockTransactionHistoryPublisher)
	mockTransactionHistoryPublisher.On("TransactionPublished", connect.Transaction{ID: "trans1"}).Return(true)

	monitor := &StationMonitor{
		Connect:                     mockConnectAPI,
		TransactionHistoryPublisher: mockTransactionHistoryPublisher,
	}

	station := connect.Station{ID: "station1"}
	monitor.transactionHistory(station)

	mockConnectAPI.AssertExpectations(t)
	mockTransactionHistoryPublisher.AssertExpectations(t)
}

func TestTransactionStatsError(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetTransactionStatistics", "station1").Return(connect.TransactionStats{}, fmt.Errorf("Error getting transaction statistics"))

	mockTransactionStatsPublisher := new(MockTransactionStatsPublisher)
	// No need to set expectations on the publisher since it should not be called

	monitor := &StationMonitor{
		Connect:                   mockConnectAPI,
		TransactionStatsPublisher: mockTransactionStatsPublisher,
	}

	station := connect.Station{ID: "station1"}
	monitor.transactionStats(station)

	mockConnectAPI.AssertExpectations(t)
	mockTransactionStatsPublisher.AssertExpectations(t)
}
func TestCreateJobsForStations(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetStations").Return([]connect.Station{{ID: "station1"}, {ID: "station2"}}, nil)

	mockTransactionStatsPublisher := new(MockTransactionStatsPublisher)
	mockTransactionHistoryPublisher := new(MockTransactionHistoryPublisher)
	mockStationStatusPublisher := new(MockStationStatusPublisher)

	ctrl := gomock.NewController(t)
	mockScheduler := gocronmocks.NewMockScheduler(ctrl)
	mockScheduler.EXPECT().NewJob(gomock.Any(), gomock.Any(), gomock.Any()).Times(4)

	monitor := &StationMonitor{
		Connect:                     mockConnectAPI,
		TransactionStatsPublisher:   mockTransactionStatsPublisher,
		TransactionHistoryPublisher: mockTransactionHistoryPublisher,
		StationStatusPublisher:      mockStationStatusPublisher,
		Scheduler:                   mockScheduler,
	}

	err := monitor.CreateJobsForStations()
	require.NoError(t, err)

	mockConnectAPI.AssertExpectations(t)
}

func TestCreateJobsForStationsError(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetStations").Return([]connect.Station{}, fmt.Errorf("Error getting stations"))

	mockTransactionStatsPublisher := new(MockTransactionStatsPublisher)
	mockTransactionHistoryPublisher := new(MockTransactionHistoryPublisher)
	mockStationStatusPublisher := new(MockStationStatusPublisher)

	ctrl := gomock.NewController(t)
	mockScheduler := gocronmocks.NewMockScheduler(ctrl)

	monitor := &StationMonitor{
		Connect:                     mockConnectAPI,
		TransactionStatsPublisher:   mockTransactionStatsPublisher,
		TransactionHistoryPublisher: mockTransactionHistoryPublisher,
		StationStatusPublisher:      mockStationStatusPublisher,
		Scheduler:                   mockScheduler,
	}

	err := monitor.CreateJobsForStations()
	require.Error(t, err)

	mockConnectAPI.AssertExpectations(t)
}

func TestMonitorStations(t *testing.T) {
	mockConnectAPI := new(MockConnectAPI)
	mockConnectAPI.On("GetStations").Return([]connect.Station{{ID: "station1"}, {ID: "station2"}}, nil)

	mockTransactionStatsPublisher := new(MockTransactionStatsPublisher)
	mockTransactionHistoryPublisher := new(MockTransactionHistoryPublisher)
	mockStationStatusPublisher := new(MockStationStatusPublisher)

	ctx, cancelCtx := context.WithCancel(context.Background())

	ctrl := gomock.NewController(t)
	mockScheduler := gocronmocks.NewMockScheduler(ctrl)
	mockScheduler.EXPECT().Start().Times(1)
	mockScheduler.EXPECT().Shutdown().Times(1).Return(nil)
	// Just accept any job creation, we're not testing that here
	mockScheduler.EXPECT().NewJob(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	monitor := &StationMonitor{
		Connect:                     mockConnectAPI,
		TransactionStatsPublisher:   mockTransactionStatsPublisher,
		TransactionHistoryPublisher: mockTransactionHistoryPublisher,
		StationStatusPublisher:      mockStationStatusPublisher,
		Scheduler:                   mockScheduler,
	}

	// Set a timer to cancel the context. Otherwise this will continue
	// indefinitely and fail when we hit the test timeout
	go func() {
		time.Sleep(1 * time.Second)
		cancelCtx()
	}()

	err := monitor.MonitorStations(ctx)

	require.EqualError(t, err, "context canceled")
	mockConnectAPI.AssertExpectations(t)
}
