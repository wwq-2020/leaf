package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type getResp struct {
	Ids []uint64 `json:"ids"`
}

func (s *service) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ids, err := s.generator.Generate(key)
	if err != nil {
		s.logger.Error("generate ids", zap.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := &getResp{Ids: ids}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.logger.Error("send get resp", zap.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
