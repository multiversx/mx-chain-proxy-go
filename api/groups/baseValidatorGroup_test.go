package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validatorPath = "/validator"

type valStatsResponseData struct {
	Statistics map[string]*data.ValidatorApiResponse `json:"statistics"`
}

// ValStatsResponse structure
type ValStatsResponse struct {
	Error string               `json:"error"`
	Data  valStatsResponseData `json:"data"`
}

func TestNewValidatorGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewValidatorGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestValidatorStatistics_ShouldErrWhenFacadeFails(t *testing.T) {
	t.Parallel()

	errStr := "expected err"
	facade := &mock.FacadeStub{
		ValidatorStatisticsHandler: func() (map[string]*data.ValidatorApiResponse, error) {
			return nil, errors.New(errStr)
		},
	}
	validatorGroup, err := groups.NewValidatorGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(validatorGroup, validatorPath)

	req, _ := http.NewRequest("GET", "/validator/statistics", nil)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.True(t, strings.Contains(response.Error, errStr))
}

func TestValidatorStatistics_ShouldWork(t *testing.T) {
	t.Parallel()

	valStatsMap := make(map[string]*data.ValidatorApiResponse)
	valStatsMap["statistics"] = &data.ValidatorApiResponse{
		NumLeaderSuccess:                   4,
		NumLeaderFailure:                   5,
		NumValidatorSuccess:                6,
		NumValidatorFailure:                7,
		NumValidatorIgnoredSignatures:      8,
		Rating:                             0.5,
		TempRating:                         0.51,
		TotalNumLeaderSuccess:              4,
		TotalNumLeaderFailure:              5,
		TotalNumValidatorSuccess:           6,
		TotalNumValidatorFailure:           7,
		TotalNumValidatorIgnoredSignatures: 8,
		ShardID:                            1,
		ValidatorStatus:                    "ok",
		RatingModifier:                     1.5,
	}
	facade := &mock.FacadeStub{
		ValidatorStatisticsHandler: func() (map[string]*data.ValidatorApiResponse, error) {
			return valStatsMap, nil
		},
	}
	validatorGroup, err := groups.NewValidatorGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(validatorGroup, validatorPath)

	req, _ := http.NewRequest("GET", "/validator/statistics", nil)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := ValStatsResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Data.Statistics["statistics"], valStatsMap["statistics"])
}
