<template>
  <div class="container">
    <br>
    <div class="row justify-content-center">
      <div class="col-12 col-md-12 col-lg-12">
        <form class="card card-sm">
          <div class="card-body row no-gutters align-items-center">
            <div class="col-auto">
              <img src="../../assets/icons/search.png" style="width: 50px; height: 50px;">
            </div>
            <div class="col">
              <input
                class="form-control form-control-lg form-control-borderless"
                type="search"
                placeholder="Enter the claim ID.. 'for example: S1'"
                v-model="searchID"
              >
            </div>
            <div class="col-auto">
              <button type="button" class="btn btn-lg btn-dark" v-on:click="searchClaims()">Search</button>
            </div>
          </div>
        </form>
      </div>
    </div>
    <div class="row" v-if="isHidden===false">
      <div class="col-12 middle-row">
        <div class="card text-left">
          <h5 class="card-header">Patient Name: {{statements.PatientName}}</h5>
          <div class="card-body">
            <table class="table table-condensed">
              <tbody>
                <tr class="tbody-font">
                  <td class="bground">Provider</td>
                  <td class="bground">{{statements.Provider}}</td>
                </tr>
                <tr class="tbody-font">
                  <td class="bground">Network</td>
                  <td class="bground">{{statements.NetworkInOut}}</td>
                </tr>
                <tr class="tbody-font">
                  <td class="bground">Payer</td>
                  <td class="bground">{{statements.PayerName}}</td>
                </tr>
                <tr class="tbody-font">
                  <td class="bground">Diagnosis</td>
                  <div class="card">
                    <table class="table table-condensed">
                    <tr class="tbody-font" v-for="diag in diagonsis" v-bind:key="diag">
                      <td class="bground">Description</td>
                      <td class="bground">{{diag.Description}}</td>
                      <td class="bground">Total Charge</td>
                      <td class="bground">{{diag.TotalCharge}}</td>
                      <td class="bground">Patient Amount</td>
                      <td class="bground">{{diag.PatientAmount}}</td>
                      <td class="bground">Billed To Insurance</td>
                      <td class="bground">{{diag.BilledToInsurance}}</td>
                    </tr>
                    </table>
                  </div>
                </tr><br/>
                <tr class="tbody-font">
                  <td class="bground">Status</td>
                  <td class="bground"><button type="submit" class="btn btn-dark">{{statements.CurrentPhase}}</button></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import axios from 'axios'

export default {
  name: 'Search',
  data () {
    return {
      statements: [],
      diagonsis: [],
      isHidden: 'false',
      searchID: null
    }
  },
  methods: {
    searchClaims: function () {
      this.isHidden = 'true'
      axios
        .get('http://129.213.162.51:2999/statements/detail/' + this.searchID)
        .then(response => {
          this.statements = response.data
          this.diagonsis = response.data.DiagnosisID
          this.isHidden = false
        })
    }
  }
}
</script>
<style scoped>
.table td, .table th {
    border: none;
}
.form-control-borderless {
  border: none;
}
.container {
  margin-top: 3%;
}
.form-control-borderless:hover,
.form-control-borderless:active,
.form-control-borderless:focus {
  border: none;
  outline: none;
  box-shadow: none;
}
.middle-row {
  margin-top: 5%;
}
.card {
  width: 100%;
}
.align {
  justify-content: left;
}
</style>
