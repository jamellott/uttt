<template>
  <div class="game">
    <div>Game ID: {{ gameID }}</div>
    <h2>
      <span :class="youClass">You</span> vs
      <span :class="oppClass">{{ getOpponent(game) }}</span>
    </h2>
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
    youClass() {
      if (this.$store.playerID == this.game.playerX) {
        return "player-x";
      }
      return "player-o";
    },
    oppClass() {
      if (this.$store.playerID == this.game.playerX) {
        return "player-o";
      }
      return "player-x";
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

h2 {
  padding: 4px;
}
h2 span {
  padding: 4px;
}
.player-o {
  background-color: lightskyblue;
}

.player-x {
  background-color: lightcoral;
}
</style>
