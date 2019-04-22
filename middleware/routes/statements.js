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
      allKeys.errorPlace = `Error fetching all keys`;
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
          answer.errorPlace = `Error fetching the statement ${allKeys[i].Key}`;
          res.status(500).send(answer);
          return;
        }
        answer = JSON.parse(answer.result)

        let diagnosisArr = [];
        for (let j = 0; j < answer.DiagnosisID.length; j++) {
          const getDiagnosis = {
            method: `POST`,
            uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/query`,
            json: true,
            body: {
              channel: `bankoforacleorderer`,
              chaincode: `healthclaims`,
              method: `queryBatchRecord`,
              args: [answer.DiagnosisID[j]]
            }
          };
          let answer2 = await rp(getDiagnosis);
          if (answer2.returnCode !== `Success`) {
            answer2.errorPlace = `Error fetching diagnosis for ${answer.DiagnosisID[i]}`;
            res.status(500).send(answer2);
            return;
          }
          diagnosisArr.push(JSON.parse(answer2.result));
        }
        answer.diagnosisArr = diagnosisArr;
        replyArr.push(answer);
      }
    }
    res.send(replyArr);
  } catch (err) {
    res.status(500).send(err);
  }
});

//// Posting of new statements.
router.post(`/`, async (req, res) => {
  // Declaring the reply object, that will hold the reply data.
  let replyObj = {};

  // Making sure all elements are there.
  if (!req.body.provider || !req.body.inNetwork || !req.body.patientName || !req.body.diagnosisIDs || !req.body.description || !req.body.currentPhase || !req.body.attachmentLink || !req.body.claimDate || !req.body.payerName) {
    replyObj.status = `Required element missing from body of request.`;
    replyObj.provider = req.body.provider || `======MISSING======`;
    replyObj.inNetwork = req.body.inNetwork || `======MISSING======`;
    replyObj.patientName = req.body.patientName || `======MISSING======`;
    replyObj.diagnosisIDs = req.body.diagnosisIDs || `======MISSING======`;
    replyObj.description = req.body.description || `======MISSING======`;
    replyObj.currentPhase = req.body.currentPhase || `======MISSING======`;
    replyObj.attachmentLink = req.body.attachmentLink || `======MISSING======`;
    replyObj.claimDate = req.body.claimDate || `======MISSING======`;
    replyObj.payerName = req.body.payerName || `======MISSING======`;
    res.status(400).send(replyObj);
    return;
  }

  const statementID = `S${Math.floor(Math.random()*100000000)}`;
  const provider = req.body.provider;
  const inNetwork = req.body.inNetwork;
  const patientName = req.body.patientName;
  let diagnosisIDs;
  if (typeof(req.body.diagnosisIDs) === `string`) {
    diagnosisIDs = [req.body.diagnosisIDs]
  } else if (typeof(req.body.diagnosisIDs) === `object`) {
    diagnosisIDs = req.body.diagnosisIDs.join(`|`);
  } else {
    replyObj.status = `diagnosisIDs type invalid.`
    res.status(400).send(replyObj);
  }
  const description = req.body.description;
  const currentPhase = req.body.currentPhase;
  const attachmentLink = req.body.attachmentLink;
  const action = `0`;
  const claimDate = req.body.claimDate;
  const payerName = req.body.payerName;

  const postStatementOptions = {
    method: `POST`,
    uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/invocation`,
    json: true,
    body: {
      channel: `bankoforacleorderer`,
      chaincode: `healthclaims`,
      method: `initStatement`,
      args: [statementID, provider, inNetwork, patientName, diagnosisIDs, description, currentPhase, attachmentLink, action, claimDate, payerName]
    }
  };

  try {
    let result = await rp(postStatementOptions);
    if (result.returnCode !== `Success`) {
      replyObj.status = `Error writing data to blockchain.`;
      replyObj.error = result;
      res.status(500).send(replyObj);
      return;
    }
    replyObj.statementID = statementID;
    res.status(200).send(replyObj);
  } catch (err) {
    replyObj.status = `Error communicating with blockchain`;
    replyObj.error = err;
    res.status(500).send(replyObj);
  }

});


/* GET a specific statment from blockchain. */
router.get(`/detail/:id`, async (req, res) => {

  const getSingleStatementOptions = {
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
    let statementInfo = await rp(getSingleStatementOptions);

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

router.patch(`/sendclaim/:id`, async (req, res) => {

  let replyObj = {};

  //First, we grab the statement to be sure it exists...
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
    statementInfo.Action = `2`;
    statementInfo.DiagnosisID = statementInfo.DiagnosisID.join(`|`);
    statementInfo.CurrentPhase = `Third Party`;

    const updateStatementOptions = {
      method: `POST`,
      uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/invocation`,
      json: true,
      body: {
        channel: `bankoforacleorderer`,
        chaincode: `healthclaims`,
        method: `updateStatement`,
        args: [statementInfo.Id, statementInfo.Provider, statementInfo.NetworkInOut, statementInfo.PatientName, statementInfo.DiagnosisID, statementInfo.Description, statementInfo.CurrentPhase, statementInfo.AttachementLink, statementInfo.Action, statementInfo.Date, statementInfo.PayerName]
      }
    };

    let result = await rp(updateStatementOptions);
    if (result.returnCode !== `Success`) {
      replyObj.status = `Error updating data on blockchain.`;
      replyObj.error = result;
      res.status(500).send(replyObj);
      return;
    }
    replyObj.status = `Statement marked as ready for 3rd party.`;
    replyObj.statementID = statementInfo.Id;
    res.status(200).send(replyObj);
  } catch (err) {
    replyObj.status = `Error occurred in promise.`;
    replyObj.error = err;
    res.status(500).send(replyObj);
  }

});


router.patch(`/adddiagnosis/:id`, async (req, res) => {

  let replyObj = {};

  if (typeof(req.body.diagnosisID) !== `object`) {
    replyObj.status = `The diagnosisID must be an array of diagnosis IDs.`;
    res.status(400).send(replyObj);
    return;
  }

  //First, we grab the statement to be sure it exists...
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
    statementInfo.DiagnosisID = statementInfo.DiagnosisID.join(`|`);
    req.body.diagnosisID.forEach((diagnosis) => {
      statementInfo.DiagnosisID += `| ${diagnosis}`;
    });

    const updateStatementOptions = {
      method: `POST`,
      uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/invocation`,
      json: true,
      body: {
        channel: `bankoforacleorderer`,
        chaincode: `healthclaims`,
        method: `updateStatement`,
        args: [statementInfo.Id, statementInfo.Provider, statementInfo.NetworkInOut, statementInfo.PatientName, statementInfo.DiagnosisID, statementInfo.Description, statementInfo.CurrentPhase, statementInfo.AttachementLink, statementInfo.Action, statementInfo.Date, statementInfo.PayerName]
      }
    };

    let result = await rp(updateStatementOptions);
    if (result.returnCode !== `Success`) {
      replyObj.status = `Error updating data on blockchain.`;
      replyObj.error = result;
      res.status(500).send(replyObj);
      return;
    }
    replyObj.status = `Diagnosis IDs added to statement.`;
    replyObj.statementID = statementInfo.Id;
    res.status(200).send(replyObj);
  } catch (err) {
    replyObj.status = `Error occurred in promise.`;
    replyObj.error = err;
    res.status(500).send(replyObj);
  }

});

router.patch(`/removediagnosis/:id`, async (req, res) => {

  let replyObj = {};

  if (typeof(req.body.diagnosisID) !== `object`) {
    replyObj.status = `The diagnosisID must be an array of diagnosis IDs.`;
    res.status(400).send(replyObj);
    return;
  }

  //First, we grab the statement to be sure it exists...
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
    
    req.body.diagnosisID.forEach((diagnosis) => {
      let theIndex = statementInfo.DiagnosisID.indexOf(diagnosis);
      if (theIndex !== -1) {
        statementInfo.DiagnosisID.splice(theIndex, 1);
      }
    });

    statementInfo.DiagnosisID = statementInfo.DiagnosisID.join(`|`);

    const updateStatementOptions = {
      method: `POST`,
      uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/invocation`,
      json: true,
      body: {
        channel: `bankoforacleorderer`,
        chaincode: `healthclaims`,
        method: `updateStatement`,
        args: [statementInfo.Id, statementInfo.Provider, statementInfo.NetworkInOut, statementInfo.PatientName, statementInfo.DiagnosisID, statementInfo.Description, statementInfo.CurrentPhase, statementInfo.AttachementLink, statementInfo.Action, statementInfo.Date, statementInfo.PayerName]
      }
    };

    let result = await rp(updateStatementOptions);
    if (result.returnCode !== `Success`) {
      replyObj.status = `Error updating data on blockchain.`;
      replyObj.error = result;
      res.status(500).send(replyObj);
      return;
    }
    replyObj.status = `Diagnosis IDs removed from statement.`;
    replyObj.statementID = statementInfo.Id;
    res.status(200).send(replyObj);
  } catch (err) {
    replyObj.status = `Error occurred in promise.`;
    replyObj.error = err;
    res.status(500).send(replyObj);
  }

});



router.patch(`/:id`, async (req, res) => {
  let replyObj = {};

  //First, we grab the statement to be sure it exists...
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

    let provider = req.body.provider || statementInfo.Provider;
    let networkInOut = req.body.inNetwork || statementInfo.NetworkInOut;
    let patientName = req.body.patientName || statementInfo.PatientName;
    let diagnosisID;
    if (req.body.diagnosisID) {
      diagnosisID = req.body.diagnosisID.join(`|`);
    } else {
      diagnosisID = statementInfo.DiagnosisID.join(`|`);
    }
    let description = req.body.description || statementInfo.Description;
    let currentPhase = req.body.currentPhase || statementInfo.CurrentPhase;
    let attachmentLink = req.body.attachmentLink || statementInfo.AttachementLink;
    let action = req.body.action || `${statementInfo.Action}`;
    let statementDate = req.body.statementDate || statementInfo.Date;
    let payerName = req.body.payerName || statementInfo.PayerName;

    const updateStatementOptions = {
      method: `POST`,
      uri: `http://132.145.136.106:3100/bcsgw/rest/v1/transaction/invocation`,
      json: true,
      body: {
        channel: `bankoforacleorderer`,
        chaincode: `healthclaims`,
        method: `updateStatement`,
        args: [statementInfo.Id, provider, networkInOut, patientName, diagnosisID, description, currentPhase, attachmentLink, action, statementDate, payerName]
      }
    };

    let result = await rp(updateStatementOptions);
    if (result.returnCode !== `Success`) {
      replyObj.status = `Error updating data on blockchain.`;
      replyObj.error = result;
      res.status(500).send(replyObj);
      return;
    }
    replyObj.status = `Statement has been updated.`;
    replyObj.statementID = statementInfo.Id;
    res.status(200).send(replyObj);
  } catch (err) {
    replyObj.status = `Error occurred in promise.`;
    replyObj.error = err;
    res.status(500).send(replyObj);
  }

});


module.exports = router;
