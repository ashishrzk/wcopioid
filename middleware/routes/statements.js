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

module.exports = router;
