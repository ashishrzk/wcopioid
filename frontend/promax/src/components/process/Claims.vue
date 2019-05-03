<template>
  <div class="container">
    <div class="row first-row">
      <hr class="line-style"/>
      <label class="title-style">CLAIMS SUMMARY REPORT</label>
      <hr class="line-style-second"/>
    </div>
    <div class="row" v-for="(statement,index) in statements" v-bind:key="statement">
      <div class="card text-left" v-if="statement.CurrentPhase === 'Third Party' && statement.NetworkInOut !== 'IN'">
        <h5 class="card-header">Payer: {{statement.PayerName}}</h5>
        <div class="card-body">
          <h5 class="card-title">Patient Name: {{statement.PatientName}}</h5>
          <table class="table table-condensed">
            <thead>
              <tr>
                <th scope="col">Patient Details</th>
                <th scope="col" class="tbody-font bground">Total Amount Billed</th>
                <th scope="col" class="tbody-font bground">Patient Amount</th>
                <th scope="col" class="tbody-font bground">Billed To Insurance</th>
                <th scope="col" class="tbody-font bground">Network</th>
                <th scope="col" class="tbody-font bground">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr class="tbody-font">
                <td>
                      <label>{{statement.Provider}}</label><br/>
                      <label>{{statement.Description}}</label>
                </td>
                <td>
                   <tr v-for="diag in statement.diagnosisArr" v-bind:key="diag" v-show = "flag == false">
                    <td><label @dblclick = "flag = true">{{diag.TotalCharge}}</label></td>
                  </tr>
                  <tr>
                    <input v-show = "flag == true" v-model="totalBilledAmount[index]">
                   </tr>
                </td>
                <td>
                  <tr v-for="diag in statement.diagnosisArr" v-bind:key="diag" v-show = "flag == false">
                    <td><label @dblclick = "flag = true">{{diag.PatientAmount}}</label></td>
                  </tr>
                  <tr>
                    <input v-show = "flag == true" v-model="patientAmount[index]">
                   </tr>
                </td>
                <td>
                  <tr v-for="diag in statement.diagnosisArr" v-bind:key="diag" v-show = "flag == false">
                    <td><label @dblclick = "flag = true">{{diag.BilledToInsurance}}</label></td>
                    </tr>
                   <tr>
                    <input v-show = "flag == true" v-model="billedtoInsurance[index]">
                   </tr>
                </td>
                <td>
                  <tr v-for="diag in statement.diagnosisArr" v-bind:key="diag">
                  <td><label>{{statement.NetworkInOut}}</label></td>
                  </tr>
                  </td>
                <td>
                  <tr v-for="diag in statement.diagnosisArr" v-bind:key="diag" v-show = "flag == false">
                    <td><label @dblclick = "flag = true">{{statement.CurrentPhase}}</label></td>
                    </tr>
                   <tr>
                    <input v-show = "flag == true" v-model="status[index]">
                   </tr>
                </td>
              </tr>
            </tbody>
          </table>
          <!-- <a href="#" class="btn btn-primary" v-on:click="sendBacktoPayer(patientAmount, billedtoInsurance, status, statement.Id, statement.DiagnosisID)">Update</a> -->
          <input type="submit" value="Update" class="btn btn-primary" v-on:click="sendBacktoPayer(statement.Id, statement.DiagnosisID, index)">
        </div>
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
      totalBilledAmount: [],
      patientAmount: [],
      billedtoInsurance: [],
      status: [],
      flag: false,
      statements: []
    }
  },
  mounted () {
    axios.get('http://129.213.162.51:2999/statements').then(response => {
      this.statements = response.data
    })
  },
  methods: {
    sendBacktoPayer (sID, dID, index) {
      console.log(this.totalBilledAmount[index])
      axios.all([
        axios.patch('http://129.213.162.51:2999/statements/' + sID, {
          currentPhase: this.status[index]
        }),
        axios.patch('http://129.213.162.51:2999/diagnosis/' + dID, {
          totalCharge: this.totalBilledAmount[index],
          patientAmount: this.patientAmount[index],
          billedToInsurance: this.billedtoInsurance[index]
        })
      ])
        .then(axios.spread((statementRes, diagnosisRes) => {
          window.location.reload()
        }))
    }
  }
}
</script>
<style scoped>
@import url("https://fonts.googleapis.com/css?family=Tangerine");
@import url("https://fonts.googleapis.com/css?family=Rajdhani");
.card {
  width: 100%;
}
.first-row {
  padding-top: 5%;
  padding-bottom: 5%;
}
.title-style {
  font-family: 'Rajdhani';
  font-size: 30px;
  font-weight: bold;
}
.table td,
.table th {
  border: none;
}
.line-style {
  width: 10%;
  margin-left: 25%;
  border-top: 10px double rgba(10, 24, 88, 0.877);
}
.line-style-second{
  width: 10%;
  border-top: 10px double rgba(10, 24, 88, 0.877);
}
</style>
