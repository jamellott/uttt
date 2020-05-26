import Vue from "vue";
import Vuex from "vuex";
import App from "./App.vue";
import ServerStore from "./serverStore.js";

Vue.use(Vuex);
Vue.config.productionTip = false;
const store = new Vuex.Store(ServerStore);

new Vue({
  render: (h) => h(App),
  store,
}).$mount("#app");
