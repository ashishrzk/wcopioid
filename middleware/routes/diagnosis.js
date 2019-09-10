var express = require(`express`);
var router = express.Router();

const rp = require(`request-promise-native`);

const authString = Buffer.from(`${process.env.RESTUSER}:${process.env.RESTPW}`).toString(`base64`);

/* GET all statements from blockchain. */
router.get(`/`, async (req, res) => {
  let replyArr = [];

  const getAllKeysOptions = {
    method: `POST`,
    uri: `https://bchost.oracle.com:3118/restproxy1/bcsgw/rest/v1/transaction/query`,
    json: true,
    body: {
      channel: `claims`,
      chaincode: `healthclaims`,
      method: `getAllRecordKeys`,
      args: []
    },
    headers: {
      Authorization : `Basic ${authString}`
    }
  };

  try {
    let allKeys = await rp(getAllKeysOptions);

    if (allKeys.returnCode !== `Success`) {
      res.status(500).send(allKeys);
      return;
    }

    allKeys = JSON.parse(allKeys.result.payload); // <- Make this a filter
  
    for (let i = 0; i < allKeys.length; i++) {
      if (allKeys[i].Key[0] === `D`) {
        const getRecord = {
          method: `POST`,
          uri: `https://bchost.oracle.com:3118/restproxy1/bcsgw/rest/v1/transaction/query`,
          json: true,
          body: {
            channel: `claims`,
            chaincode: `healthclaims`,
            method: `queryBatchRecord`,
            args: [allKeys[i].Key]
          },
          headers: {
            Authorization : `Basic ${authString}`
          }
        };
        let answer = await rp(getRecord);
        if (answer.returnCode !== `Success`) {
          res.status(500).send(answer);
          return;
        }
        replyArr.push(JSON.parse(answer.result.payload));
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
  if (!req.body.CPTCode || !req.body.CPTName || !req.body.ICD9 || !req.body.description || !req.body.totalCharge || !req.body.patientAmount || !req.body.billedToInsurance) {
    replyObj.status = `Required element missing from body of request.`;
    replyObj.CPTCode = req.body.CPTCode || `======MISSING======`;
    replyObj.CPTName = req.body.patientName || `======MISSING======`;
    replyObj.ICD9 = req.body.diagnosisIDs || `======MISSING======`;
    replyObj.description = req.body.description || `======MISSING======`;
    replyObj.totalCharge = req.body.currentPhase || `======MISSING======`;
    replyObj.patientAmount = req.body.attachmentLink || `======MISSING======`;
    replyObj.billedToInsurance = req.body.claimDate || `======MISSING======`;
    res.status(400).send(replyObj);
    return;
  }

  const diagnosisID = `D${Math.floor(Math.random()*100000000)}`;
  const CPTCode = req.body.CPTCode;
  const CPTName = req.body.CPTName;
  const ICD9 = req.body.ICD9;
  const description = req.body.description;
  const totalCharge = req.body.totalCharge;
  const patientAmount = req.body.patientAmount;
  const billedToInsurance = req.body.billedToInsurance;

  const postDiagnosisOptions = {
    method: `POST`,
    uri: `https://bchost.oracle.com:3118/restproxy1/bcsgw/rest/v1/transaction/invocation`,
    json: true,
    body: {
      channel: `claims`,
      chaincode: `healthclaims`,
      method: `initDiagnosis`,
      args: [diagnosisID, CPTCode, CPTName, ICD9, description, totalCharge, patientAmount, billedToInsurance]
    },
    headers: {
      Authorization : `Basic ${authString}`
    }
  };

  try {
    let result = await rp(postDiagnosisOptions);
    if (result.returnCode !== `Success`) {
      replyObj.status = `Error writing data to blockchain.`;
      replyObj.error = result;
      res.status(500).send(replyObj);
      return;
    }
    replyObj.diagnosisID = diagnosisID;
    res.status(200).send(replyObj);
  } catch (err) {
    replyObj.status = `Error communicating with blockchain`;
    replyObj.error = err;
    res.status(500).send(replyObj);
  }

});


router.patch(`/:id`, async (req, res) => {
  let replyObj = {};

  //First, we grab the statement to be sure it exists...
  const getDiagnosisOptions = {
    method: `POST`,
    uri: `https://bchost.oracle.com:3118/restproxy1/bcsgw/rest/v1/transaction/query`,
    json: true,
    body: {
      channel: `bankoforacleorderer`,
      chaincode: `healthclaims`,
      method: `queryBatchRecord`,
      args: [req.params.id]
    },
    headers: {
      Authorization : `Basic ${authString}`
    }
  };

  try {
    let diagnosisInfo = await rp(getDiagnosisOptions);

    if (diagnosisInfo.returnCode !== `Success`) {
      res.status(500).send(diagnosisInfo);
      return;
    }

    diagnosisInfo = JSON.parse(diagnosisInfo.result.payload);

    let CPTCode = req.body.CPTCode || diagnosisInfo.CPTCode;
    let CPTName = req.body.CPTName || diagnosisInfo.CPTName;
    let ICD9 = req.body.ICD9 || diagnosisInfo.ICD9;
    let description = req.body.description || diagnosisInfo.Description;
    let totalCharge = req.body.totalCharge || diagnosisInfo.TotalCharge;
    let patientAmount = req.body.patientAmount || diagnosisInfo.PatientAmount;
    let billedToInsurance = req.body.billedToInsurance || diagnosisInfo.BilledToInsurance;

    const updateDiagnosisOptions = {
      method: `POST`,
      uri: `https://bchost.oracle.com:3118/restproxy1/bcsgw/rest/v1/transaction/invocation`,
      json: true,
      body: {
        channel: `claims`,
        chaincode: `healthclaims`,
        method: `updateDiagnosis`,
        args: [diagnosisInfo.Id, CPTCode, CPTName, ICD9, description, totalCharge, patientAmount, billedToInsurance]
      },
      headers: {
        Authorization : `Basic ${authString}`
      }
    };

    let result = await rp(updateDiagnosisOptions);
    if (result.returnCode !== `Success`) {
      replyObj.status = `Error updating data on blockchain.`;
      replyObj.error = result;
      res.status(500).send(replyObj);
      return;
    }
    replyObj.status = `Diagnosis has been updated.`;
    replyObj.diagnosisID = diagnosisInfo.Id;
    res.status(200).send(replyObj);
  } catch (err) {
    replyObj.status = `Error occurred in promise.`;
    replyObj.error = err;
    res.status(500).send(replyObj);
  }

});





module.exports = router;
