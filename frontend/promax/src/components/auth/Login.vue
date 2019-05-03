<template>
  <section>
    <div class="login-bg-top">
      <div>
        <img class="d-block img-fluid" src="../../assets/promax-befunk.jpg">
      </div>
      <div class="logo-circle"></div>
      <div class="container">
        <div class="row">
          <div class="col-md-12">
            <h3 class="title">Promax</h3>
          </div>
        </div>
      </div>
    </div>
    <hr class="line-style">
    <div class="container">
      <div class="col-md-12">
        <div class="bg-content">
          <form class="form-style">
            <div class="form-group row">
              <label for="staticEmail" class="col-sm-2 col-form-label">USERNAME:</label>
              <div class="col-sm-10">
                <input type="text" v-model="input.username" name="username" class="form-control" id="staticEmail">
              </div>
            </div>
            <div class="form-group row">
              <label for="inputPassword" class="col-sm-2 col-form-label">PASSWORD:</label>
              <div class="col-sm-10">
                <input
                  type="password"
                  class="form-control"
                  id="inputPassword"
                  name="password" v-model="input.password"
                >
              </div>
            </div>
            <br>
            <div>
              <button type="submit" class="btn btn-primary" v-on:click="login()">SIGN IN</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </section>
</template>
<script>
export default {
  name: 'Login',
  data () {
    return {
      input: {
        username: '',
        password: ''
      }
    }
  },
  methods: {
    login: function () {
      if (this.input.username === 'test' && this.input.password === 'test') {
        this.$router.push({name: 'Home'})
      } else {
        console.log('Invalid username and password.')
      }
    }
  }
}
</script>
<style scoped>
@import url("https://fonts.googleapis.com/css?family=Assistant");
@import url("https://fonts.googleapis.com/css?family=Monoton");
.login-bg {
  height: 110vh;
  background-color: #f5f5f5;
}
.login-bg-top {
  margin-top: -5%;
}
.logo-circle {
  background-image: url("../../assets/promax-p-edited.png");
  background-repeat: no-repeat;
  background-size: cover;
  border-radius: 125px;
  width: 150px;
  height: 150px;
  margin-top: -4.8%;
  display: inline-block;
  border: 5px solid rgba(211, 211, 211, 0.9);
  animation-name: bounceIn;
  animation-duration: 2s;
}
.title {
  font-family: "Monoton";
  font-size: 75px;
  color: rgba(10, 24, 88, 0.877);
  letter-spacing: 10px;
  animation-name: bounceIn;
  animation-duration: 2s;
}
.bg-content {
  text-align: center;
  padding-top: 50px;
  padding-left: 110px;
  padding-right: 110px;
  margin-top: 40px;
  margin-bottom: 90px;
  background-color: #fff;
  min-height: 300px;
}
.line-style {
  width: 10%;
  margin-left: 45%;
  border-top: 5px dotted rgba(10, 24, 88, 0.877);
}
.form-style {
  font-family: "Assistant";
  font-weight: bold;
  font-size: 14px;
  width: 48%;
  margin-left: 25%;
}

@keyframes bounceIn {
  0% {
    transform: scale(0.1);
    opacity: 0;
  }
  60% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(1);
  }
}
</style>
