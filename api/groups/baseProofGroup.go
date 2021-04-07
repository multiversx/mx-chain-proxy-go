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
		{Path: "/verify", Handler: pg.verifyProof, Method: http.MethodPost},
	}
	pg.baseGroup.endpoints = baseRoutesHandlers

	return pg, nil
}

func (pg *proofGroup) getProof(c *gin.Context) {
	rootHash := c.Param("roothash")
	address := c.Param("address")

	proof, err := pg.facade.GetProof([]byte(rootHash), []byte(address))
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"proof": proof}, "", data.ReturnCodeSuccess)
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

	ok, err := pg.facade.VerifyProof(proofParams.RootHash, proofParams.Address, proofParams.Proof)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"ok": ok}, "", data.ReturnCodeSuccess)
}
