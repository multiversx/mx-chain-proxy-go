package groups

import (
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type proofGroup struct {
	facade ProofFacadeHandler
	*baseGroup
}

func NewProofGroup(facadeHandler data.FacadeHandler) (*proofGroup, error) {
	facade, ok := facadeHandler.(ProofFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	pg := &proofGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/root-hash/:roothash/address/:address", Handler: pg.getProof, Method: http.MethodGet},
		{Path: "/address/:address", Handler: pg.getProofCurrentRootHash, Method: http.MethodGet},
		{Path: "/verify", Handler: pg.verifyProof, Method: http.MethodPost},
	}
	pg.baseGroup.endpoints = baseRoutesHandlers

	return pg, nil
}

func (pg *proofGroup) getProof(c *gin.Context) {
	rootHash := c.Param("roothash")
	if rootHash == "" {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrEmptyRootHash.Error(), data.ReturnCodeRequestError)
		return
	}

	address := c.Param("address")
	if address == "" {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrEmptyAddress.Error(), data.ReturnCodeRequestError)
		return
	}

	getProofResp, err := pg.facade.GetProof(rootHash, address)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, getProofResp)
}

func (pg *proofGroup) getProofCurrentRootHash(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrEmptyAddress.Error(), data.ReturnCodeRequestError)
		return
	}

	getProofResp, err := pg.facade.GetProofCurrentRootHash(address)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, getProofResp)
}

func (pg *proofGroup) verifyProof(c *gin.Context) {
	proofParams := &data.VerifyProofRequest{}
	err := c.ShouldBindJSON(proofParams)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
		)
		return
	}

	verifyProofResp, err := pg.facade.VerifyProof(proofParams.RootHash, proofParams.Address, proofParams.Proof)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, verifyProofResp)
}
