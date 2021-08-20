package bnc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type StakingReward struct {
	Total         int64          `json:"total"`
	RewardDetails []RewardDetail `json:"rewardDetails"`
}

type RewardDetail struct {
	ChainId    string  `json:"chainId"`
	Validator  string  `json:"validator"`
	ValName    string  `json:"valName"`
	Delegator  string  `json:"delegator"`
	Reward     float64 `json:"reward"`
	Height     int64   `json:"height"`
	RewardTime string  `json:"rewardTime"`
}

func GetStakingReward(url string) (*StakingReward, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sr := new(StakingReward)
	if err := json.Unmarshal(body, sr); err != nil {
		return nil, err
	}

	return sr, nil
}
