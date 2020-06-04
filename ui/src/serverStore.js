class WSMessage {
  constructor(type, payload) {
    this.messageType = type;
    this.payload = payload;
  }
}

function createLoginRequestVerifier(socket, resolve, reject) {
  let handler = (ev) => {
    socket.removeEventListener("message", handler);
    console.log(ev);
    let data = JSON.parse(ev.data);
    if (data.messageType === "LoginSuccess") {
      resolve(socket);
    } else if (data.messageType === "LoginFailure") {
      console.error(data.payload.message);
      reject("login failed");
    } else {
      reject("Unexpected MessageType" + data.messageType);
    }
  };

  socket.addEventListener("message", handler);
}

function sendLoginRequest(socket, username) {
  let msg = new WSMessage("LoginRequest", {
    playerID: username,
  });
  socket.send(JSON.stringify(msg));
  return new Promise((resolve, reject) => {
    createLoginRequestVerifier(socket, resolve, reject);
  });
}

const webSocketHandler = {
  socket: null,
  store: null,
  installFunc() {
    let handler = this;
    return (store) => {
      handler.store = store;
    };
  },
  setSocket(socket) {
    if (this.socket !== null) {
      // todo: support reconnections? Or just refresh
      console.error("socket being replaced!");
    }
    this.socket = socket;
    let handler = this;
    socket.addEventListener("message", (ev) => {
      handler.handleMessage(ev.data);
    });
  },
  handleMessage(msg) {
    switch (msg.type) {
      case "GameUpdate":
        this.state.commit("playMove", msg);
        break;
      default:
        console.error("unknown websocket message type: " + msg.type);
        break;
    }
  },
  sendMessage(msg) {
    this.socket.send(msg);
  },
};

const store = {
  state: {
    username: "",
    userID: "",
    games: [],
  },
  mutations: {
    setUser(state, { username, userID }) {
      state.username = username;
      state.userID = userID;
    },
    initGammes(state, games) {
      state.games = games;
    },
    gameUpdate(state, game) {
      let idx = state.games.find((g) => g.id == game.id);
      state.games[idx] = game;
    },
  },
  plugins: [webSocketHandler.installFunc()],
  actions: {
    login(context, username) {
      let url = process.env.VUE_APP_SERVER_SOCKET_URL;
      console.log("connecting to " + url);
      let socket = new WebSocket(url);
      return new Promise((resolve, reject) => {
        socket.addEventListener("open", () => {
          sendLoginRequest(socket, username)
            .then((socket) => {
              console.log("successful connection");
              webSocketHandler.setSocket(socket);
              resolve();
            })
            .catch((reason) => {
              console.error(reason);
              reject(reason);
            });
        });
      });
    },
    playMove(context, move) {
      let message = new WSMessage("PlayMove", move);
      webSocketHandler.sendMessage(message);
    },
  },
};

export default store;
