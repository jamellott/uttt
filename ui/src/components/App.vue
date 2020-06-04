<template>
  <div class="login">
    <div class="row">
      <div class="col-12 col-3-md">
        <ul>
          <li class="row"><profile-display /></li>
          <li class="row" v-for="game in games" :key="game.id">
            <game-list-item />
          </li>
          <li class="row"><new-game-button /></li>
        </ul>
      </div>
      <div class="col-9-md col-12">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script>
import GameListItem from "./sidebar/GameListItem.vue";
import NewGameButton from "./sidebar/NewGameButton.vue";
import ProfileDisplay from "./sidebar/ProfileDisplay.vue";

export default {
  name: "App",
  components: {
    GameListItem,
    NewGameButton,
    ProfileDisplay,
  },
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
  border-right: 2px #bbb;
}

li {
  padding: 20px;
}
li div {
  display: inline-block;
  margin: 10px;
  border: 4px #ddd;
  border-radius: 10px;
  background-color: #555;
  color: #ccc;
  min-height: 100px;
  width: 100%;
}
</style>
