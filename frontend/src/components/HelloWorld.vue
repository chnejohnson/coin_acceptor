<template>
  <h1>{{ msg }}</h1>
  <p style="font-size: 50px">{{ amount }} 元</p>
</template>

<script>
import { ref } from "vue";

export default {
  name: "HelloWorld",
  props: {
    msg: String,
  },
  setup() {
    const amount = ref(0);

    var ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = function (e) {
      console.log(e.data);
      amount.value = e.data;
    };

    return { amount };
  },
};
</script>
