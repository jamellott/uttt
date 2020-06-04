<template>
  <div class="login">
    <div class="container">
      <div class="row">
        <div class="col-12 col-3-md">
          <ul>
            <li><profile-display /></li>
            <li v-for="game in games" :key="game.id"><game-list-item /></li>
            <li><new-game-button /></li>
          </ul>
        </div>
        <div class="col-9-md col-12">
          <router-view></router-view>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "App",
  data() {
    return {
      username: "",
    };
  },
  methods: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
