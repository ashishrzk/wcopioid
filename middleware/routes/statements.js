var express = require(`express`);
var router = express.Router();

const rp = require(`request-promise-native`);

/* GET all statements from blockchain. */
router.get(`/`, async (req, res) => {
  let replyArr = [];

  const getAllKeysOptions = {
    method: `POST`,
    uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/query`,
    json: true,
    body: {
      channel: `bankoforacleorderer`,
      chaincode: `healthclaims`,
      method: `getAllRecordKeys`,
      args: []
    }
  };

  try {
    let allKeys = await rp(getAllKeysOptions);

    if (allKeys.returnCode !== `Success`) {
      res.status(500).send(allKeys);
      return;
    }

    allKeys = JSON.parse(allKeys.result);
  
    for (let i = 0; i < allKeys.length; i++) {
      if (allKeys[i].Key[0] === `S`) {
        const getRecord = {
          method: `POST`,
          uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/query`,
          json: true,
          body: {
            channel: `bankoforacleorderer`,
            chaincode: `healthclaims`,
            method: `queryBatchRecord`,
            args: [allKeys[i].Key]
          }
        };
        let answer = await rp(getRecord);
        if (answer.returnCode !== `Success`) {
          res.status(500).send(answer);
          return;
        }
        replyArr.push(JSON.parse(answer.result));
      }
    }
    res.send(replyArr);
  } catch (err) {
    res.status(500).send(err);
  }
});

//// TODO: Make this work.
router.post(`/`, async (req, res) => {
  const statementID = `S${Math.floor(Math.random()*100000000)}`;
  const provider = req.params.provider;
  const inNetwork = req.params.inNetwork;
  const patientName = req.params.patientName;
  const diagnosisIDs = req.params.diagnosisIDs;
  const description = req.params.description;
  const currentPhase = req.params.currentPhase;
  const attachmentLink = req.params.attachmentLink;
  const action = req.params.action;
  const claimDate = req.params.claimDate;
});

/* GET all statements from blockchain. */
router.get(`/detail/:id`, async (req, res) => {

  const getAllKeysOptions = {
    method: `POST`,
    uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/query`,
    json: true,
    body: {
      channel: `bankoforacleorderer`,
      chaincode: `healthclaims`,
      method: `queryBatchRecord`,
      args: [req.params.id]
    }
  };

  try {
    let statementInfo = await rp(getAllKeysOptions);

    if (statementInfo.returnCode !== `Success`) {
      res.status(500).send(statementInfo);
      return;
    }

    statementInfo = JSON.parse(statementInfo.result);
    let diagnosisArr = [];
    for (let i = 0; i < statementInfo.DiagnosisID.length; i++) {
      const getDiagnosis = {
        method: `POST`,
        uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/query`,
        json: true,
        body: {
          channel: `bankoforacleorderer`,
          chaincode: `healthclaims`,
          method: `queryBatchRecord`,
          args: [statementInfo.DiagnosisID[i]]
        }
      };
      let answer = await rp(getDiagnosis);
      if (answer.returnCode !== `Success`) {
        res.status(500).send(answer);
        return;
      }
      diagnosisArr.push(JSON.parse(answer.result));
    }

    statementInfo.DiagnosisID = diagnosisArr;
    res.send(statementInfo);
  } catch (err) {
    res.status(500).send(err);
  }
});




module.exports = router;
