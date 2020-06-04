<template>
  <div class="container">
    <div
      v-for="(gridRow, i) in game.grids"
      :key="i"
      :class="`row grid-row p-1 game-row-${i + 1}`"
    >
      <div
        v-for="(gridCol, i) in gridRow"
        :key="i"
        :class="`col-4 grid-col game-col-${i + 1}`"
      >
        <div
          v-for="(row, i) in gridCol.squares"
          :key="i"
          :class="`row sub-row game-row-${i + 1}`"
        >
          <div
            v-for="(sq, i) in row"
            :key="i"
            :class="`col-4 p-0 game-col-${i + 1}`"
          >
            <square :square="sq" :game="game" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Square from "./Square.vue";

export default {
  name: "Grid",
  props: {
    game: Object,
  },
  components: {
    Square,
  },
  data() {
    return {
      username: "",
    };
  },
  computed: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
/*
.row {
  border-radius: 2px;
}
*/

.grid-row {
  min-height: 150px;
}

.game-row-1 {
  border-bottom: 3px solid black;
}

.game-row-2 {
  border-bottom: 3px solid black;
}

.game-row-3 {
}

.sub-row {
  min-height: 40px;
}
.game-col-1 {
  border-right: 3px solid black;
}

.game-col-2 {
  border-right: 3px solid black;
}

.game-col-3 {
}
</style>
