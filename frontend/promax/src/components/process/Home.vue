<template>
  <section>
    <div>
      <img class="d-block img-fluid" src="../../assets/home_edited.jpg" style="height: auto;">
    </div>
    <div class="row">
      <div class="col-12">
        <label class="title-style hvr-grow">Welcome to Promax</label>
      </div>
    </div>
    <div class="row third-row">
        <div class="col-md-1 col-1-margin-left"></div>
      <div class="col-md-2 hvr-grow" v-on:click="claims()">
        <img src="../../assets/icons/claims.png" style="width: 100px; height: 100px;">
        <br>Claims
      </div>
      <div class="col-md-1">
        <label class="vl"></label>
      </div>
      <div class="col-md-2 hvr-grow">
        <img src="../../assets/icons/providers.png" style="width: 100px; height: 100px;">
        <br>Providers
      </div>
      <div class="col-md-1">
        <label class="vl"></label>
      </div>
      <div class="col-md-2 hvr-grow">
        <img src="../../assets/icons/payers.png" style="width: 100px; height: 100px;">
        <br>Payers
      </div>
      <div class="col-md-1">
        <label class="vl"></label>
      </div>
      <div class="col-md-2 hvr-grow">
        <img src="../../assets/icons/help.png" style="width: 100px; height: 100px;">
        <br>Help
      </div>
    </div>
  </section>
</template>
<script>
export default {
  name: 'Home',
  methods: {
    claims: function () {
      this.$router.push({name: 'Claims'})
    }
  }
}
</script>
<style scoped>
@import url("https://fonts.googleapis.com/css?family=Tangerine");
@import url("https://fonts.googleapis.com/css?family=Voltaire");
.vl {
  border-left: 2px solid rgba(8, 8, 8, 0.466);
  height: 100px;
}
.title-style {
    padding: 3%;
    font-family: 'Tangerine';
    font-size: 90px;
}
.col-1-margin-left {
    margin-left: -4%;
}
.third-row {
    font-weight: bold;
    font-family: 'Voltaire'
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
</style>
