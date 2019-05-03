<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-12"><label class="welcome hvr-grow">Hello, there!</label></div>
    </div>
    <div class="row middle-row">
      <div class="col-md-3" v-on:click="search()">
          <div class="col-md-12 hvr-bounce-to-top">
              <img src="../../assets/icons/search.png" style="width: 100px; height: 100px;"/><br/>
              <label>Search</label>
              </div>
      </div>
      <div class="col-md-3" v-on:click="claims()">
          <div class="col-md-12 hvr-bounce-to-top">
              <img src="../../assets/icons/claims.png" style="width: 100px; height: 100px;"/><br/>
              <label>Claims</label>
              </div>
      </div>
      <div class="col-md-3">
          <div class="col-md-12 hvr-bounce-to-top">
              <img src="../../assets/icons/provider.png" style="width: 100px; height: 100px;"/><br/>
              <label>Providers</label>
              </div>
      </div>
    </div>
    <div class="row bottom-row">
      <div class="col-md-3">
          <div class="col-md-12 hvr-bounce-to-top">
              <img src="../../assets/icons/pharmacy.png" style="width: 100px; height: 100px;"/><br/>
              <label>Pharmacies</label>
              </div>
      </div>
      <div class="col-md-3">
          <div class="col-md-12 hvr-bounce-to-top">
              <img src="../../assets/icons/resources.png" style="width: 100px; height: 100px;"/><br/>
              <label>Resources</label>
              </div>
      </div>
      <div class="col-md-3">
          <div class="col-md-12 hvr-bounce-to-top">
              <img src="../../assets/icons/help.png" style="width: 100px; height: 100px;"/><br/>
              <label>Help</label>
          </div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: 'Home',
  methods: {
    search: function () {
      this.$router.push({name: 'Search'})
    },
    claims: function () {
      this.$router.push({name: 'Claims'})
    }
  }
}
</script>
<style scoped>
.container-fluid {
  padding: 0 !important
}

.welcome {
    font-size: 35px;
    padding: 4%;
    letter-spacing: 5px;
    font-weight: bold;
    color: #2098D1;
}
.middle-row {
    margin-left: 20%;
}
.bottom-row {
    margin-left: 20%;
    margin-top: 1.7%;
}
.col-md-12 {
    padding: 20%;
    background-color: white;
    border-radius: 10px;
    font-weight: bold;
    color: #2098D1;
}
.hvr-grow {
    display: inline-block;
    vertical-align: middle;
    transform: translateZ(0);
    box-shadow: 0 0 1px rgba(0, 0, 0, 0);
    backface-visibility: hidden;
    -moz-osx-font-smoothing: grayscale;
    transition-duration: 0.3s;
    transition-property: transform;
}

.hvr-grow:hover,
.hvr-grow:focus,
.hvr-grow:active {
    transform: scale(1.1);
}
.hvr-bounce-to-top {
  display: inline-block;
  vertical-align: middle;
  -webkit-transform: perspective(1px) translateZ(0);
  transform: perspective(1px) translateZ(0);
  box-shadow: 0 0 1px rgba(0, 0, 0, 0);
  position: relative;
  -webkit-transition-property: color;
  transition-property: color;
  -webkit-transition-duration: 0.5s;
  transition-duration: 0.5s;
}
.hvr-bounce-to-top:before {
  content: "";
  position: absolute;
  z-index: -1;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #2098D1;
  -webkit-transform: scaleY(0);
  transform: scaleY(0);
  -webkit-transform-origin: 50% 100%;
  transform-origin: 50% 100%;
  -webkit-transition-property: transform;
  transition-property: transform;
  -webkit-transition-duration: 0.5s;
  transition-duration: 0.5s;
  -webkit-transition-timing-function: ease-out;
  transition-timing-function: ease-out;
}
.hvr-bounce-to-top:hover, .hvr-bounce-to-top:focus, .hvr-bounce-to-top:active {
  color: white;
}
.hvr-bounce-to-top:hover:before, .hvr-bounce-to-top:focus:before, .hvr-bounce-to-top:active:before {
  -webkit-transform: scaleY(1);
  transform: scaleY(1);
  -webkit-transition-timing-function: cubic-bezier(0.52, 1.64, 0.37, 0.66);
  transition-timing-function: cubic-bezier(0.52, 1.64, 0.37, 0.66);
}
</style>
