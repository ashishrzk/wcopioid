<template>
  
    <div style="position: absolute; width:100%; height:100%" >
        <v-card>
          <h2 style="font-weight: bold;">Patient Risk Profile Analysis</h2>
        <oracle-dv project-path="/shared/Drug Abuse/Drug_Abuse_V5" active-page="canvas" active-tab-id="1">
        </oracle-dv>
        <h2 style="font-weight: bold;">Patient Prescription Analysis</h2>
        <oracle-dv project-path="/shared/Drug Abuse/Drug Abuse Use Case 2 Final" active-page="canvas" active-tab-id="1">
        </oracle-dv>
        </v-card>
        
    </div>
</template>


<script>
module.exports = {
  data: function () {
    return {
      greeting: 'Hello'
    }
  },
  mounted() {
    let recaptchaScript = document.createElement('script')
    recaptchaScript.setAttribute('src', 'http://150.136.203.67/dv/ui/api/v1/plugins/embedding/standalone/embedding.js')
    document.head.appendChild(recaptchaScript)
    window.addEventListener('load', function(event) {
     requirejs(['knockout', 'ojs/ojcore', 'ojs/ojknockout', 'ojs/ojcomposite', 'jet-composites/oracle-dv/loader'], function(ko) {   ko.applyBindings();
     });
   });
  },

}
</script>



<style scoped>
p {
  font-size: 2em;
  text-align: center;
}
</style>