class WSMessage {
  constructor(type, payload) {
    this.messageType = type;
    this.payload = payload;
  }
}

const MessageType = {
  LoginRequest: "LoginRequest",
  LoginSuccess: "LoginSuccess",
  LoginFailure: "ErrorMessage",
  ErrorMessage: "ErrorMessage",
};

function createLoginRequestVerifier(socket, resolve, reject) {
  let handler = (ev) => {
    socket.removeEventListener("message", handler);
    console.log(ev);
    let data = JSON.parse(ev.data);
    if (data.messageType === MessageType.LoginSuccess) {
      resolve(socket);
    } else if (data.messageType === MessageType.LoginFailure) {
      console.error(data.payload.message);
      reject("login failed");
    } else {
      reject("Unexpected MessageType: " + data.messageType);
    }
  };

  socket.addEventListener("message", handler);
}

function sendLoginRequest(socket, username) {
  let msg = new WSMessage(MessageType.LoginRequest, {
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
      default:
        console.error("unknown websocket message type: " + msg.type);
        break;
    }
  },
};

const store = {
  state: {},
  mutations: {},
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
  },
};

export default store;
