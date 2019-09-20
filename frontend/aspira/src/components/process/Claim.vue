<template>
  <div class="container-fluid" style="background: rgb(255, 255, 255);">
    <div class="row first-row">
      <label>Patient Claim: {{claim.CLAIM_NUMBER}}</label>
    </div>
        <div class="row second-row" style="float: right; margin:0 30px 0 0;">
          <v-card>
              <ul>
                <h3 v-if="claim.DOSAGE_BASE_LINE_FLAG || claim.TOTAL_RISK_FACTOR=='High Risk' || claim.TOTAL_RISK_FACTOR=='Moderate Risk'">This claim has been flagged because:     </h3>
                <div v-if="claim.TOTAL_RISK_FACTOR=='High Risk'">- The patient has been identified as a high-risk profile.</div>
                <div v-else-if="claim.TOTAL_RISK_FACTOR=='Moderate Risk'">- The patient has been identified as a moderate-risk profile.</div>
                <div style="padding-bottom: 10px" v-if="claim.DOSAGE_BASE_LINE_FLAG">- The prescribed dose is above the normal threshold.</div>
                <button type="button" class="btn btn-danger" v-on:click="viewanalytics()">View Analytics</button>
              </ul>
            </v-card>

        </div>
        <div class="row second-row">
          <span class="border border-dark" style="padding: 15px">
          <span class="claim-table">
          <label>
          <span style="font-weight: bold">Name: </span>{{claim.NAME}}
          </label>
          <label style="padding-left: 40px">
          <span style="font-weight: bold">Gender: </span>{{claim.GENDER}}
          </label>
          <br>
          <label>
          <span style="font-weight: bold">Home Phone: </span>919-425-4523
          </label>
          <br>
          <label>
            <span style="font-weight: bold">Address: </span>{{claim.ADDRESS1}}
          </label>
          <br>
          <label>
          <span style="font-weight: bold">City: </span>{{claim.CITY}}
          </label>
          <label style="padding-left: 40px">
          <span style="font-weight: bold">State: </span>{{claim.STATE}}
          </label>
          <label style="padding-left: 40px">
          <span style="font-weight: bold">Zip: </span>{{claim.ZIP}}
          </label>
          <br>
          <label>
          <span style="font-weight: bold">Date of Injury: </span>{{claim.EVENTDATE}}
          </label>
          <br>
          <label>
          <span style="font-weight: bold">Employer: </span>{{claim.EMPLOYER}}
          </label>
          <br>
          <label>
          <span style="font-weight: bold">Employer Phone: </span>849-458-3425
          </label>
          <br>
          <label>
          <span style="font-weight: bold">Nature: </span>{{claim.NATURETITLE}}
          </label>
          <br>
          <label>
          <span style="font-weight: bold">Part of Body: </span>{{claim.PART_OF_BODY_TITLE}}
          </label>
          </span>
          </span>
        </div>
        <div class="row second-row" style="padding-right: 150px;">
        <span class="border border-dark" style="padding: 20px">
          <v-card>
            <label>
            <h2>Investigation Report </h2>
            </label>
            <br>
            <label>
            <span style="font-weight: bold">Investigator Name: </span>Dr. Pepper
            </label>
            <label style="padding-left: 40px">
            <span style="font-weight: bold">Investigation Date: </span>9/14/2019
            </label>
            <br><br>
            <label for="detailsofinc" style="font-weight: bold">Details of the incident </label><br>
            <div>{{claim.FINAL_NARRATIVE}}</div>
            <br>
            <label for="emphist" style="font-weight: bold">Employment History </label><br>
            <div>Current Employer: {{claim.EMPLOYER}}</div>
            <br>
            <label for="medicalhistory" style="font-weight: bold">Medical History </label><br>
            <div>Substance Abuse: {{claim.PREVIOUS_SUBSTANCE_ABUSE_HISTORY}}</div>
            <div>Mental Disorders: {{claim.PREVIOUS_MENTAL_HISTORY}}</div>
            <div>Family History of Abuse: {{claim.FAMILY_HISTORY_OF_ABUSE}}</div>
            <br>
            <label for="socialfactors" style="font-weight: bold">Social Factors </label><br>
            <div>Marital Status: {{claim.MARITAL_STATUS}}</div>
            <div>Dependants: {{claim.DEPENDANTS}}</div>
            <br>
          </v-card>
        </span>
        </div>
        <div>
          <button type="button" class="btn btn-danger">REJECT CLAIM</button>
          <button type="button" class="btn btn-dark" @click="showModal">SEND FOR REVIEW</button>
          <button type="button" class="btn btn-success">APPROVE CLAIM</button>
        </div>
        <br><br><br>
  </diV>
</template>

<style scoped>
.first-row {
  padding-top: 1%;
  padding-bottom: 1%;
  padding-left: 10%;
  background: rgba(0, 0, 0, 0.089);
  font-size: 32px;
  letter-spacing: 2px;
  font-weight: bold;
}
.second-row {
  padding-top: .5%;
  padding-bottom: .5%;
  padding-left: 11%;
  background: rgb(255, 255, 255);
  font-size: 16px;
  text-align: left;
}
.flagged {
  background: rgb(255, 255, 255);
  font-size: 16px;
  text-align: left
}
.claim-table {
  font-size: 20px;
}
</style>


<script>
import axios from 'axios'
import modal from '@/components/process/modal.vue'

export default {
  name: 'Claim',
  components: {
      modal,
    },
  data: function () {
  //console.dir(this.$route.params['claim'])
    return {
      isModalVisible: false,
      claim: [{"CLAIM_NUMBER": "SC937225"}],
      //claim: this.$route.params['claim'],
      route1: this.$route.params.claim_number
      
    }
  },
  mounted () {
    //console.log(this.$route)
    let url = 'https://cors-anywhere.herokuapp.com/http://129.213.147.141:3050/riskandclaims/'+this.$route.params.claim_number
    console.log(url)
    axios.get(url).then(response => {
      this.claim = response.data[0]
      console.log(this.claim);
    })
  },
  methods:{
    viewanalytics(){
      //this.$router.push({name: 'Home'})
      this.$router.push({name: 'PAnalytics'})
    },
    showModal() {
        this.isModalVisible = true;
      },
      closeModal() {
        this.isModalVisible = false;
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
  padding-top: 2%;
  padding-bottom: 2%;
  padding-left: 11%;
  background: rgba(0, 0, 0, 0.089);
  font-size: 16px;
  text-align: center;
}

</style>
