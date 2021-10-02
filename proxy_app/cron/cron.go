package cron

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"

	proxytypes "github.com/unibrightio/proxy-api/types"

	common "github.com/unibrightio/proxy-api/common"
	"github.com/unibrightio/proxy-api/dbutil"
	"github.com/unibrightio/proxy-api/logger"

	businesslogic "github.com/unibrightio/proxy-api/business_logic"
)

func queryTrustmeshes() {
	logger.Info("query trustmesh start")

	var trustmeshEntries []proxytypes.TrustmeshEntry
	dbutil.Db.GetConn().Where("commitment_state=?", common.UncommittedCommitmentState).Find(&trustmeshEntries)

	logger.Infof("found %v trustmesh entries\n", len(trustmeshEntries))
	var jobs = make(chan proxytypes.Job, len(trustmeshEntries))
	var results = make(chan proxytypes.Result, len(trustmeshEntries))
	createWorkerPool(1, jobs, results)

	for _, trustmeshEntry := range trustmeshEntries {
		logger.Infof("creating job for %v\n", trustmeshEntry.TransactionHash)
		job := proxytypes.Job{TrustmeshEntry: trustmeshEntry}
		jobs <- job
	}
	close(jobs)

	for result := range results {
		logger.Infof("Tx hash %v, height %v, timestamp %v\n", result.Job.TrustmeshEntry.TransactionHash, result.TxInfo.TxHeight, result.TxInfo.TxTimestamp)
	}

	logger.Info("query trustmesh end")
}

func getTxInfo(txHash string) (txInfo *proxytypes.TxInfo, err error) {
	// fetching tx details
	str := "http://" + viper.Get("TENDERMINT_API_URL").(string) + "/tx?hash=0x" + txHash
	httpRes, err := http.Get(str)
	if err != nil {
		logger.Errorf("error during http tx req %v\n", err)
		return &proxytypes.TxInfo{}, errors.New("get tx info error")
	}

	// if transaction is not found it is not yet committed
	if httpRes.StatusCode != 200 {
		logger.Error("tx not committed yet")
		return &proxytypes.TxInfo{TxCommitted: false}, errors.New("get tx info error")
	}

	// if it's found should be committed at this point, decode
	var committedTx proxytypes.TxResp
	err = json.NewDecoder(httpRes.Body).Decode(&committedTx)
	if err != nil {
		logger.Errorf("error decoding tx %v\n", err)
		return &proxytypes.TxInfo{}, errors.New("error decoding tx")
	}
	// query for block at specific height to find timestamp
	str = "http://" + viper.Get("TENDERMINT_API_URL").(string) + "/block?height" + committedTx.TxResult.Height
	httpRes, err = http.Get(str)
	if err != nil {
		logger.Errorf("error during http block req %v\n", err)
		return &proxytypes.TxInfo{}, errors.New("get blcok info error")
	}
	var commitedBlock proxytypes.BlockResp
	err = json.NewDecoder(httpRes.Body).Decode(&commitedBlock)
	if err != nil {
		logger.Errorf("error decoding block %v\n", err)
		return &proxytypes.TxInfo{}, errors.New("error decoding block")
	}
	logger.Infof("DECODED COMMITTED TX HEIGHT %v AND TIMESTAMP %v\n", committedTx.TxResult.Height, commitedBlock.BlockResult.Block.Header.Time)
	txValid := true
	if committedTx.TxResult.TxResultInfo.Code != 0 {
		txValid = false
	}
	return &proxytypes.TxInfo{
		TxHeight:    committedTx.TxResult.Height,
		TxTimestamp: commitedBlock.BlockResult.Block.Header.Time,
		TxCommitted: true,
		TxValid:     txValid,
		TxCode:      committedTx.TxResult.TxResultInfo.Code,
		TxLog:       committedTx.TxResult.TxResultInfo.Log,
	}, nil
}

func worker(jobs chan proxytypes.Job, results chan proxytypes.Result) {
	defer close(results)
	for job := range jobs {
		txInfo, err := getTxInfo(job.TrustmeshEntry.TransactionHash)
		if err != nil {
			// here it would be http error
			// it seems that we can just let it go through result channel
			logger.Warnf("result error %v\n", err)
		}
		logger.Infof("result tx %v transaction type %v\n", txInfo, job.TrustmeshEntry.BaseledgerTransactionType)
		output := proxytypes.Result{Job: job, TxInfo: *txInfo}
		businesslogic.ExecuteBusinessLogic(output)
		results <- output
	}
}

func createWorkerPool(noOfWorkers int, jobs chan proxytypes.Job, results chan proxytypes.Result) {
	for i := 0; i < noOfWorkers; i++ {
		go worker(jobs, results)
	}
}

func StartCron() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Seconds().SingletonMode().Do(queryTrustmeshes)

	s.StartAsync()
}
