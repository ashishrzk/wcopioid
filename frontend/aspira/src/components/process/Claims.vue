<template>
  <div class="container-fluid">
    <div class="row first-row">
      <label>Patient Claims: Brief Summary</label>
    </div>
    <div class="row third-row">
      <div style="width: 80%; padding-bottom: 20px" v-for="claim in this.claims" v-bind:key="claim">
        <div @click="clickMethod(claim)" class="border">
          <table style="width: 100%">
            <thead>
            <tr>
              <th scope="col" style="width: 40%">Claim Details</th>
              <th scope="col" class="tbody-font bground" style="width: 15%">Date of Injury</th>
              <th scope="col" class="tbody-font bground" style="width: 15%">Nature</th>
              <th scope="col" class="tbody-font bground" style="width: 15%">Part of Body</th>
              <th scope="col" class="tbody-font bground" style="width: 15%">Flagged</th>
            </tr>
            </thead>
            <tbody>
              <tr class="tbody-font">
                <th scope="row">
                  <label>{{claim.claim}} {{claim.CLAIM_NUMBER}}</label>
                  <ul class="detail-style">
                    <li><b>Employer:</b> {{claim.EMPLOYER}}</li>
                    <li><b>Gender:</b> {{claim.GENDER}}</li>
                    <li><b>Location:</b> {{claim.STATE}}</li>
                  </ul>
                </th>
                <td class="bground">
                {{claim.EVENTDATE}}
                </td>
                <td class="bground">
                {{claim.NATURETITLE}}
                </td>
                <td class="bground">
                {{claim.PART_OF_BODY_TITLE}}
                </td>
                <td class="bground">
                  <button title="View Details" v-if="claim.DOSAGE_BASE_LINE_FLAG || claim.TOTAL_RISK_FACTOR=='High Risk' || claim.TOTAL_RISK_FACTOR=='Moderate Risk'" type="button" class="btn btn-danger" v-on:click="clickMethod(statement.Id)">Yes</button>
                  <button v-else type="button" class="btn btn-success" v-on:click="clickMethod(statement.Id)">No</button>
                </td>
              </tr>
            </tbody>
          </table>
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
      claims: [],
      claimID: null,
    }
  },
  mounted () {
    axios.get('https://cors-anywhere.herokuapp.com/http://129.213.147.141:3050/riskandclaims').then(response => {
      this.claims = response.data
      console.log(this.claims);
    })
  },
  methods: {
    clickMethod (claim) {
       //add code that you wish to happen on click
       //console.log("clicked is" + claim.CLAIM_NUMBER)
       //this.$router.push({name: 'Claim', params: { claim } }) //oG pass in all data
       let claim_number = claim.CLAIM_NUMBER
       //this.$router.push({name: 'claims', params: { claim_number } })
       this.$router.push({ path: `/claims/${claim_number}`}) //works
       //this.$router.push({ path: '/claims/', query: { id: 'private' } })
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
  padding-top: 2%;
  padding-right: 11%;
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
  padding-left: 8%;
}
.bground {
  background-color: rgba(0, 0, 0, 0.11);
}
</style>
