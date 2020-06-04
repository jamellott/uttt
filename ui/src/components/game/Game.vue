<template>
  <div class="game">
    <div>Game {{ gameID }} vs {{ getOpponent(game) }}</div>
    <div>
      <grid :game="game"></grid>
    </div>
  </div>
</template>

<script>
import Grid from "./Grid.vue";

export default {
  name: "Game",
  props: ["gameID"],
  components: {
    Grid,
  },
  data() {
    return {
      username: "",
    };
  },
  computed: {
    game() {
      return this.$store.state.games[this.gameID];
    },
  },
  methods: {
    getOpponent(game) {
      if (game.playerX == this.$store.state.playerID) {
        return game.playerOName;
      }

      return game.playerXName;
    },
  },
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
