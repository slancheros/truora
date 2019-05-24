<template>
  <div id="app">
    <DomainInfo msg="Welcome to your Domain Info App"/>
  </div>
</template>

<script>
  import BootstrapVue from 'bootstrap-vue'
  import Vue from "vue"
  import axios from "axios"
  import DomainInfo from './components/DomainInfo.vue'
  Vue.use(BootstrapVue)
  Vue.use(axios)
  import 'bootstrap/dist/css/bootstrap.css'
  import 'bootstrap-vue/dist/bootstrap-vue.css'

export default {
  name: 'app',
  components: {
    DomainInfo
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
