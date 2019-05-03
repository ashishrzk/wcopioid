<template>
  <div class="container-fluid">
    <div class="row first-row">
      <label>Patient Claims: Brief Summary</label>
    </div>
    <div class="row second-row">
      <ul style="min-width: 100px;">
        <li>Find the patients claims and the amount charged for services.</li>
        <li>Calculate the amount paid by patient for in network services.</li>
        <li>Send the claims in batch for out of network services.</li>
        <li>Generate report for individual patients.</li>
      </ul>
    </div>
    <div class="row third-row">
      <div class="col-11" v-for="statement in statements" v-bind:key="statement">
        <table v-if="statement.CurrentPhase !== 'Outside Review'" class="table table-condensed">
          <thead>
            <tr>
              <th scope="col" style="width: 45%">Patient Details</th>
              <th scope="col" class="tbody-font bground">Network</th>
              <th scope="col" class="tbody-font bground">Current Phase</th>
              <th scope="col" class="tbody-font bground">Total Amount Billed</th>
              <th scope="col" class="tbody-font bground">Patient Amount</th>
              <th scope="col" class="tbody-font bground">Billed To Insurance</th>
              <th scope="col" class="tbody-font bground">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr class="tbody-font">
              <th scope="row">
                <label>{{statement.Id}} {{statement.PatientName}}</label>
                <ul class="detail-style">
                  <li>
                    <b>Provider:</b>
                    {{statement.Provider}}
                  </li>
                  <li>
                    <b>Description:</b>
                    {{statement.Description}}
                  </li>
                  <li>
                    <b>Attachment:</b>
                    <a>{{statement.AttachementLink}}</a>
                  </li>
                </ul>
              </th>
              <td class="bground">
                <tr><td>{{statement.NetworkInOut}}</td></tr>
                </td>
              <td class="bground">
                <tr><td>{{statement.CurrentPhase}}</td></tr>
                </td>
              <td class="bground">
                  <tr v-for="diag in statement.diagnosisArr" v-bind:key="diag">
                    <td>{{diag.TotalCharge}}</td>
                  </tr>
              </td>
              <td class="bground">
                  <tr v-for="diag in statement.diagnosisArr" v-bind:key="diag">
                    <td>{{diag.PatientAmount}}</td>
                  </tr>
              </td>
              <td class="bground">
                  <tr v-for="diag in statement.diagnosisArr" v-bind:key="diag">
                    <td>{{diag.BilledToInsurance}}</td>
                  </tr>
              </td>
              <td class="bground">
                <button v-if="statement.CurrentPhase === 'Ready to Claim'" type="button" class="btn btn-danger" v-on:click="sendToThirdParty(statement.Id)">Send Claim</button>
                <button v-else-if="statement.CurrentPhase === 'Third Party'" type="button" class="btn btn-warning">Processing</button>
                <button v-else type="button" class="btn btn-success">Success</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
<script>
import axios from 'axios'

export default {
  name: 'Claims',
  data () {
    return {
      statements: [],
      statementID: null
    }
  },
  mounted () {
    axios.get('http://129.213.162.51:2999/statements').then(response => {
      this.statements = response.data
    })
  },
  methods: {
    sendToThirdParty (statementID) {
      axios.patch('http://129.213.162.51:2999/statements/sendclaim/' + statementID).then(response => {
        console.log('Sent ' + statementID + ' for third party processing successfully.')
        window.location.reload()
      })
    }
  }
}
</script>
<style scoped>
.first-row {
  padding-top: 1%;
  padding-bottom: 1%;
  padding-left: 10%;
  background: rgba(0, 0, 0, 0.089);
  font-size: 24px;
  letter-spacing: 2px;
  font-weight: bold;
}
.second-row {
  padding-top: 2%;
  padding-bottom: 2%;
  padding-left: 11%;
  background: rgb(255, 255, 255);
  font-size: 16px;
  text-align: left;
}
.third-row {
  padding-bottom: 2%;
  padding-left: 11%;
  background: rgb(255, 255, 255);
  font-size: 16px;
  text-align: left;
}
.col-8 {
  min-width: 60%;
  font-size: 13px;
}
.col-4 {
  font-size: 13px;
}
.card-header {
  font-weight: bold;
}
.tbody-font {
  font-size: 14px;
}
.table td,
.table th {
  border: none;
}
table.table.table-condensed {
  border: 1px solid rgba(0, 0, 0, 0.301);
}
.detail-style {
  font-weight: lighter;
  padding-left: 5%;
}
.bground {
  background-color: rgba(0, 0, 0, 0.11);
}
</style>
