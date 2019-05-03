<template>
<div class="full-pg">
    <div class="container">
    <div class="row inner-row-style">
      <div class="col-lg-8 img-padding">
          <img class="d-block img-fluid" src="../../assets/healthcare.jpg">
      </div>
      <div class="col-lg-4 form-margin">
        <form>
          <div class="form-group">
            <div>
              <label class="login-title">Aspira Insurance</label>
            </div>
            <br>
            <div>
              <label>USERNAME:</label>
              <input type="text" name="username" v-model="input.username">
            </div>
            <div>
              <label>PASSWORD:</label>
              <input type="password" name="password" v-model="input.password">
            </div>
            <br>
            <div>
              <button type="submit" class="btn btn-dark" v-on:click="login()">SIGN IN</button>
            </div>
          </div>
        </form>
      </div>
      </div>
      </div>
  </div>
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
@import url('https://fonts.googleapis.com/css?family=Tangerine');
@import url('https://fonts.googleapis.com/css?family=Assistant');
.container {
  animation-name: bounceIn;
  animation-duration: 2s;
}
.navbar {
    display: none !important;
}
.full-pg {
  background-image: url('../../assets/healthcare-blur.jpg');
  background-repeat: no-repeat;
  background-size: cover;
  background-position: center;
  height: 101vh;
  /* margin-top: -1% !important; */
  overflow: auto;
}
.row {
  background-color: rgb(255, 255, 255);
  border-radius: 10px;
  margin-top: 20%;
  width: 100%;
}
.radius {
  border-radius: 1%;
}
.img-padding {
  padding: 2%;
}
.img-fluid {
  border-radius: 10px;
}
.login-title {
  font-family: 'Tangerine';
  font-size: 60px;
  font-weight: bold;
  color: rgb(112, 157, 182)
}
.form-margin {
  margin-top: 12%;
  font-family: 'Assistant';
  font-weight: bold;
  font-size: 14px;
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
