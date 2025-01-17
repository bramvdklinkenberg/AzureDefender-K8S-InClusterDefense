package cache

import (
	"context"
	"github.com/Azure/AzureDefender-K8S-InClusterDefense/pkg/infra/instrumentation"
	"github.com/Azure/AzureDefender-K8S-InClusterDefense/pkg/infra/retrypolicy"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

// We'll be able to store suite-wide
// variables and add methods to this
// test suite struct
type TestSuiteRedisCache struct {
	suite.Suite
}

const (
	_key   = "hello"
	_value = "world"
)

var (
	_ctx             = context.Background()
	_retryPolicy     retrypolicy.IRetryPolicy
	_client          *RedisCacheClient
	_redisClientMock *redis.Client
	_redisMock       redismock.ClientMock
)

func (suite *TestSuiteRedisCache) SetupTest() {
	// TODO retry mocking in all places
	_redisClientMock, _redisMock = redismock.NewClientMock()
	retryPolicyConfiguration := &retrypolicy.RetryPolicyConfiguration{RetryAttempts: 1, RetryDurationInMS: 10}
	_retryPolicy = retrypolicy.NewRetryPolicy(instrumentation.NewNoOpInstrumentationProvider(), retryPolicyConfiguration)
	_client = NewRedisCacheClient(instrumentation.NewNoOpInstrumentationProvider(), _redisClientMock, _retryPolicy)
}

func (suite *TestSuiteRedisCache) Test_Get_KeyIsExist_ShouldReturnValue() {
	// Setup
	expectedValue := _value

	_redisMock.ExpectGet(_key).SetVal(expectedValue)

	// Act
	actual, err := _client.Get(_ctx, _key)

	// Test
	suite.Nil(err)
	suite.Equal(expectedValue, actual)
}

func (suite *TestSuiteRedisCache) Test_Get_KeyIsNotExist_ShouldReturnErr() {
	// Setup
	_redisMock.ExpectGet(_key).SetErr(redis.Nil)

	// Act
	_, err := _client.Get(_ctx, _key)

	// Test
	suite.NotNil(err)
}

func (suite *TestSuiteRedisCache) Test_Set_NewKey_ShouldReturnNil() {
	// Setup
	duration := time.Duration(3)
	_redisMock.ExpectSet(_key, _value, duration).RedisNil()

	// Act
	err := _client.Set(_ctx, _key, _value, duration)
	suite.Nil(err)
}

func (suite *TestSuiteRedisCache) Test_Set_NegativeExpiration_ShouldReturnErr() {
	// Setup
	duration := time.Duration(-3)
	_redisMock.ExpectSet(_key, _value, duration).SetVal(_value)

	// Act
	err := _client.Set(_ctx, _key, _value, duration)
	suite.IsType(&NegativeExpirationCacheError{}, err)
}

func (suite *TestSuiteRedisCache) assertExpectationsMocks() {
	//_retryPolicyMock.AssertExpectations(suite.T())

}

// We need this function to kick off the test suite, otherwise
// "go test" won't know about our tests
func TestRedisCacheClient(t *testing.T) {
	suite.Run(t, new(TestSuiteRedisCache))
}
