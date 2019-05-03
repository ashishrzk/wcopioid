<template>
  <div id="app">
    <nav class="navbar navbar-expand-lg">
      <a class="navbar-brand nav-title" v-if="$route.name != 'Login'">PROMAX</a>
      <button
        class="navbar-toggler"
        type="button"
        data-toggle="collapse"
        data-target="#navbarTogglerDemo02"
        aria-controls="navbarTogglerDemo02"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse">
        <ul class="navbar-nav ml-auto">
          <li class="nav-item">
            <router-link :to="{name: 'Home'}" class="nav-link"><img src="../../assets/icons/home.png" style="width: 40px; height: 40px;"/></router-link>
          </li>
          <li class="nav-item">
            <router-link :to="{name: 'Login'}" class="nav-link"><img src="../../assets/icons/logout.png" style="width: 40px; height: 40px;"/></router-link>
          </li>
        </ul>
      </div>
    </nav>
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
@import url("https://fonts.googleapis.com/css?family=Merienda");
.navbar {
  background: transparent;
}
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
body {
  background-color: white;
}
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: 100vh;
}
.flex-center {
  display: flex;
  justify-content: center;
}
.nav-title {
    font-family: 'Merienda';
    font-weight: bold;
    font-size: 24px;
    letter-spacing: 8px;
    color: #2098D1;
}
.navbar {
    background-color: white;
    border-bottom: 3px solid rgba(0, 0, 0, 0.233);
}
</style>
